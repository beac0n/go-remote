package server

import (
	"crypto/cipher"
	"encoding/binary"
	"go-remote/src/util"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func Run(port string, keyBase64 string, timeFrame int64, quitChan chan bool) {
	util.InitConfigDir()
	util.InitTimestampFile()

	keyFileBytes := util.GetKeyBytes(keyBase64)
	aesKeyBytes := keyFileBytes[0:util.AesKeySize]
	expectedSourcePort := strconv.Itoa(util.GetSourcePort(aesKeyBytes))

	aeadKey, err := util.GetAesGcmAEAD(aesKeyBytes)
	util.Check(err, "could not parse aes key bytes")

	packetConnection := setupPacketConnection(port)
	defer packetConnection.Close()
	for {
		if quit(quitChan) {
			return
		}

		encryptedBytes := make([]byte, util.EncryptedDataLen+1)
		n, address, err := packetConnection.ReadFrom(encryptedBytes)
		if err != nil {
			log.Println("ERROR could not read packet from connection:", err)
			continue
		}

		log.Println("# incoming:", address)

		addressSplit := strings.Split(address.String(), ":")
		sourcePort := addressSplit[len(addressSplit)-1]

		if sourcePort != expectedSourcePort {
			log.Println("ERROR expected source port", expectedSourcePort, "but got", sourcePort)
			continue
		}

		if n != util.EncryptedDataLen {
			log.Println("ERROR received invalid bytes length. Expected", util.EncryptedDataLen, "got", n)
			continue
		}

		if validateIncomingData(encryptedBytes[0:util.EncryptedDataLen], aeadKey, timeFrame) {
			_, err := net.Dial("unix", util.SocketPath)
			if err != nil {
				log.Println("ERROR: could not connect to unix socket: ", err)
			}

			emptyBuffer(packetConnection)
		}
	}
}

func quit(quit chan bool) bool {
	if quit == nil {
		return false
	}
	select {
	case <-quit:
		return true
	default:
		return false
	}
}

func setupPacketConnection(port string) net.PacketConn {

	var packetConnection net.PacketConn
	var err error

	if os.Getenv("LISTEN_PID") == strconv.Itoa(os.Getpid()) {
		log.Println("Starting UDP server on socket from systemd")
		f := os.NewFile(3, "from systemd")
		packetConnection, err = net.FilePacketConn(f)
	} else {
		address := "0.0.0.0:" + port
		log.Println("Starting UDP server on", address)
		packetConnection, err = net.ListenPacket("udp", address)
	}

	util.Check(err, "could not start udp server")

	return packetConnection
}

func validateIncomingData(encryptedBytes []byte, aeadKey cipher.AEAD, timeFrame int64) bool {
	dataBytes, err := util.DecryptData(aeadKey, encryptedBytes)
	if err != nil {
		return false
	}

	lastTsNanoInt := getLastTsNanoInt()
	nowNanoInt := time.Now().UnixNano()

	if lastTsNanoInt >= nowNanoInt {
		log.Fatal("ERROR: last timestamp must be smaller than now")
		return false
	}

	tsNanoBytes := dataBytes[0:util.TimestampLen]
	tsNanoInt := int64(binary.LittleEndian.Uint64(tsNanoBytes))

	timeFrameNanoInt := timeFrame * util.SecInNs
	startTsNano := nowNanoInt - timeFrameNanoInt
	endTsNano := nowNanoInt + timeFrameNanoInt

	isWithinTimeFrame := startTsNano < tsNanoInt && endTsNano > tsNanoInt
	if isWithinTimeFrame && tsNanoInt > lastTsNanoInt {
		util.WriteTimestampFile(tsNanoBytes)
		return true
	} else if !isWithinTimeFrame {
		log.Println("ERROR timestamp not within timeframe.\nExpected",
			tsNanoInt, "(", time.Unix(tsNanoInt/util.SecInNs, 0), ")\nto be between",
			startTsNano, "(", time.Unix(startTsNano/util.SecInNs, 0), ")\nand",
			endTsNano, "(", time.Unix(endTsNano/util.SecInNs, 0), ")")
		return false
	} else {
		log.Println("ERROR got invalid timestamp. Expected", tsNanoInt,
			"(", time.Unix(tsNanoInt/util.SecInNs, 0), ") to be greater than the last timestamp")
		return false
	}
}

func getLastTsNanoInt() int64 {
	lastTimestamp, err := util.ReadTimestampFile()
	util.Check(err, "could not read timestamp file")

	return int64(binary.LittleEndian.Uint64(lastTimestamp))
}

func emptyBuffer(con net.PacketConn) {
	err := con.SetReadDeadline(time.Now().Add(time.Second))
	if err != nil {
		log.Println("ERROR could not SetReadDeadline:", err)
		return
	}

	n := 1
	buffer := make([]byte, util.EncryptedDataLen)

	for n > 0 {
		n, _, err = con.ReadFrom(buffer)
		if err != nil {
			break
		}
	}

	err = con.SetReadDeadline(time.Time{})
	if err != nil {
		log.Println("ERROR could not SetReadDeadline to zero:", err)
	}
}
