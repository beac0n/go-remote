package server

import (
	"bytes"
	"crypto/ed25519"
	"encoding/binary"
	"go-remote/src/util"
	"io/ioutil"
	"log"
	"net"
	"os/exec"
	"strings"
	"time"
)

func Run(port *string, keyFilePath *string, timeFrame *int64, command *string, timeout *int64, end *string) {
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
			executeCommand(command)
			time.Sleep(time.Duration(*timeout) * time.Second)
			executeCommand(end)
			emptyBuffer(packetConnection, util.EncryptedDataLen)
		}
	}
}

func emptyBuffer(con net.PacketConn, bytesLength int) {
	n := 1
	buffer := make([]byte, bytesLength)
	for n > 0 {
		n, _, _ = con.ReadFrom(buffer)
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

	timestampBytes := dataBytes[0:util.TimestampLen]
	timestampInt := int64(binary.LittleEndian.Uint64(timestampBytes))

	isValid := isTsWithinTimeFrame(timestampInt, timeFrame) && isCurrentTsGreaterLastTs(timestampInt)
	if isValid {
		util.WriteBytes(util.FilePathTimestamp, timestampBytes)
	}

	return isValid
}

func isTsWithinTimeFrame(timestampInt int64, timeFrame *int64) bool {
	now := time.Now().UnixNano()
	timeframeNanoSeconds := *timeFrame * 1000000000
	return now-timeframeNanoSeconds < timestampInt && now > timestampInt
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
