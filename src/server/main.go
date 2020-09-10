package main

import (
	"bytes"
	"crypto/ed25519"
	"encoding/binary"
	"flag"
	"go-remote/src/util"
	"io/ioutil"
	"log"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func main() {
	keyFilePath := flag.String("key", "", "path to key file")
	port := flag.String("port", "8080", "udp port")
	timeFrame := flag.Int64("timeframe", int64(5), "timestamp in request must not be older than this timeframe (in seconds)")
	commandStart := flag.String("command-start", "echo 'start!'", "the command to execute before the command timeout")
	commandTimeout := flag.Int64("command-timeout", int64(60), "how long to wait before executing the end command")
	commandEnd := flag.String("command-end", "echo 'end!'", "the command to execute after the command timeout")

	flag.Parse()

	run(port, keyFilePath, timeFrame, commandStart, commandTimeout, commandEnd)
}

func run(port *string, keyFilePath *string, timeFrame *int64, commandStart *string, commandTimeout *int64, commandEnd *string) {
	serverKeyFileBytes := util.ReadBytes(*keyFilePath)

	serverKeyBytes := serverKeyFileBytes[0:util.ServerKeyLen]
	cryptoKeyBytes := serverKeyFileBytes[util.ServerKeyLen:util.ServerKeyFileLen]

	packetConnection := setupPacketConnection(port)

	for {
		encryptedBytes := make([]byte, util.EncryptedDataLen)
		n, address, err := packetConnection.ReadFrom(encryptedBytes)
		if err != nil {
			log.Println("ERROR could not read packet from connection:", err)
			continue
		}

		log.Println(address.String() + ": ")

		if n != util.EncryptedDataLen {
			log.Println("ERROR received incorrect bytes length")
			continue
		}

		if validateIncomingData(encryptedBytes, serverKeyBytes, cryptoKeyBytes, timeFrame) {
			executeCommand(commandStart)
			time.Sleep(time.Duration(*commandTimeout) * time.Second)
			executeCommand(commandEnd)
			emptyBuffer(packetConnection, util.EncryptedDataLen)
		}
	}
}

func emptyBuffer(con net.PacketConn, bytesLength int) {
	n := 1
	buffer := make([]byte, bytesLength)

	err := con.SetReadDeadline(time.Now().Add(time.Second))
	if err != nil {
		log.Println("ERROR could not SetReadDeadline:", err)
		return
	}

	for n > 0 {
		n, _, err = con.ReadFrom(buffer)
		if err != nil {
			log.Println("WARN could not ReadFrom connection:", err)
			break
		}
	}

	err = con.SetReadDeadline(time.Time{})
	if err != nil {
		log.Println("ERROR could not SetReadDeadline to zero:", err)
	}
}

func setupPacketConnection(port *string) net.PacketConn {
	address := "0.0.0.0:" + *port
	log.Println("Starting UDP server on " + address)

	packetConnection, err := net.ListenPacket("udp", address)
	util.Check(err, "could not start udp server")

	return packetConnection
}

func validateIncomingData(encryptedBytes []byte, serverKey []byte, cryptoKeyBytes []byte, timeFrame *int64) bool {
	dataBytes, success := util.DecryptData(cryptoKeyBytes, encryptedBytes)
	if !success {
		return false
	}

	messageBytes, signatureBytes := dataBytes[:util.MsgLen], dataBytes[util.MsgLen:util.DataLen]
	if !ed25519.Verify(serverKey, messageBytes, signatureBytes) {
		log.Println("ERROR got invalid signature")
		return false
	}

	tsNanoBytes := dataBytes[0:util.TimestampLen]
	tsNanoInt := int64(binary.LittleEndian.Uint64(tsNanoBytes))
	tsNanoStr := strconv.FormatInt(tsNanoInt, 10)

	nowNanoInt := time.Now().UnixNano()
	nowNanoStr := strconv.FormatInt(nowNanoInt, 10)

	timeframeNanoSeconds := *timeFrame * util.SecInNs
	startTsNano := nowNanoInt - timeframeNanoSeconds
	startTsNanoStr := strconv.FormatInt(startTsNano, 10)

	withinTimeFrame := startTsNano < tsNanoInt && nowNanoInt > tsNanoInt
	currentTsGreaterLastTs := isCurrentTsGreaterLastTs(tsNanoInt)

	isValid := withinTimeFrame && currentTsGreaterLastTs
	if isValid {
		util.WriteBytes(util.FilePathTimestamp, tsNanoBytes)
	} else if !withinTimeFrame {
		log.Println("ERROR got invalid timestamp.\n" +
			"Expected " + tsNanoStr + " (" + time.Unix(tsNanoInt/util.SecInNs, 0).String() + ")\n" +
			"to be between " + startTsNanoStr + " (" + time.Unix(startTsNano/util.SecInNs, 0).String() + ")\n" +
			"and " + nowNanoStr + " (" + time.Unix(nowNanoInt/util.SecInNs, 0).String() + ")")
	} else if !currentTsGreaterLastTs {
		log.Println("ERROR got invalid timestamp. " +
			"Expected " + tsNanoStr + " to be greater than the last timestamp")
	}

	return isValid
}

func isCurrentTsGreaterLastTs(timestampInt int64) bool {
	isCurrentTsGreaterLastTs := true
	lastTimestamp, err := ioutil.ReadFile(util.FilePathTimestamp)
	if err == nil {
		lastTimestampInt := int64(binary.LittleEndian.Uint64(lastTimestamp))
		isCurrentTsGreaterLastTs = timestampInt > lastTimestampInt
	} else {
		log.Println("WARN: could not read timestamp file:", err)
	}
	return isCurrentTsGreaterLastTs
}

func executeCommand(command *string) {
	commandSplit := strings.Split(*command, " ")
	commandSplitLen := len(commandSplit)

	var cmd *exec.Cmd
	if commandSplitLen == 0 {
		return
	} else if commandSplitLen == 1 {
		cmd = exec.Command(commandSplit[0])
	} else {
		cmd = exec.Command(commandSplit[0], commandSplit[1:]...)
	}

	var outBytes bytes.Buffer
	cmd.Stdout = &outBytes

	log.Println("running command " + *command)
	err := cmd.Run()
	if err == nil {
		log.Println(outBytes.String())
	} else {
		log.Println("ERROR when running command:", err)
	}
}
