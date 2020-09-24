package server

import (
	"bytes"
	"crypto/cipher"
	"encoding/binary"
	"go-remote/src/util"
	"log"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func Run(port string, keyFilePath string, timeFrame int64, commandStart string, commandTimeout int64, commandEnd string, quitChan chan bool) {
	initTimestampFile()

	keyFileBytes := util.ReadBytes(keyFilePath)
	aesKeyBytes := keyFileBytes[0:util.AesKeySize]
	expectedSourcePort := strconv.Itoa(util.GetSourcePort(aesKeyBytes))

	aeadKey, err := util.GetAesGcmAEAD(aesKeyBytes)
	util.Check(err, "could not parse aes key bytes")

	packetConnection := setupPacketConnection(port)
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
			executeCommand(commandStart)
			time.Sleep(time.Duration(commandTimeout) * time.Second)
			executeCommand(commandEnd)
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

func initTimestampFile() {
	_, err := util.ReadTimestampFile()
	if err != nil {
		util.WriteTimestampFile(util.GetTimestampNowBytes())
	}

}

func setupPacketConnection(port string) net.PacketConn {
	address := "0.0.0.0:" + port
	log.Println("Starting UDP server on", address)

	packetConnection, err := net.ListenPacket("udp", address)
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
		log.Println("ERROR got invalid timestamp.\nExpected",
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

func executeCommand(command string) {
	commandSplit := strings.Split(command, " ")
	commandSplitLen := len(commandSplit)

	var cmd *exec.Cmd
	if commandSplitLen == 0 {
		return
	} else if commandSplitLen == 1 {
		cmd = exec.Command(commandSplit[0])
	} else {
		cmd = exec.Command(commandSplit[0], commandSplit[1:]...)
	}

	var stdOutBytes bytes.Buffer
	var stdErrBytes bytes.Buffer
	cmd.Stdout = &stdOutBytes
	cmd.Stderr = &stdErrBytes

	log.Println("running command", command)
	if err := cmd.Run(); err != nil {
		log.Println("ERROR when running command:", err)
	}
	log.Println("Stdout:", stdOutBytes.String())
	log.Println("Stderr:", stdErrBytes.String())
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
