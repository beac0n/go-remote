package server

import (
	"bytes"
	"crypto/ed25519"
	"encoding/binary"
	"go-remote/src/util"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func Run(port *string, keyId *string, timeFrame *int64, command *string) {
	serverKeyFileBytes := util.GetFileBytes(*keyId, util.ServerSuffix)

	serverKeyBytes := serverKeyFileBytes[0:util.ServerKeyLen]
	cryptoKeyBytes := serverKeyFileBytes[util.ServerKeyLen:util.ServerKeyFileLen]

	packetConnection := setupPacketConnection(port)

	for {
		if validateIncomingData(packetConnection, serverKeyBytes, cryptoKeyBytes, timeFrame) {
			executeCommand(command)
		}
	}

}

func setupPacketConnection(port *string) net.PacketConn {
	address := "0.0.0.0:" + *port
	log.Println("Starting UDP server on " + address)

	packetConnection, err := net.ListenPacket("udp", address)
	util.Check(err, "could not start udp server")

	return packetConnection
}

func validateIncomingData(packetConnection net.PacketConn, serverKey []byte, cryptoKeyBytes []byte, timeFrame *int64) bool {
	encryptedBytes := make([]byte, util.EncryptedDataLen)
	n, address, err := packetConnection.ReadFrom(encryptedBytes)

	if err != nil {
		log.Println("ERROR could not read packet from connection:", err)
		return false
	}

	addressString := address.String()
	if n != util.EncryptedDataLen {
		log.Println("ERROR received " + strconv.Itoa(n) + " bytes from " + addressString + ", but expected " +
			strconv.Itoa(util.EncryptedDataLen) + " bytes")
		return false
	}

	dataBytes, success := util.DecryptData(cryptoKeyBytes, encryptedBytes)
	if !success {
		return false
	}

	messageBytes, signatureBytes := dataBytes[:util.MsgLen], dataBytes[util.MsgLen:util.DataLen]
	if !ed25519.Verify(serverKey, messageBytes, signatureBytes) {
		log.Println("ERROR got invalid signature from " + addressString)
		return false
	}

	timestampBytes := dataBytes[:util.TimestampLen]
	timestampInt := int64(binary.LittleEndian.Uint64(timestampBytes))

	isValid := isTsWithinTimeFrame(timestampInt, timeFrame) && isCurrentTsGreaterLastTs(timestampInt)
	if isValid {
		updateTimestampFile(timestampBytes)
	}

	return isValid
}

func isTsWithinTimeFrame(timestampInt int64, timeFrame *int64) bool {
	now := time.Now().UnixNano()
	return now-*timeFrame < timestampInt && now > timestampInt
}

func updateTimestampFile(timestampBytes []byte) {
	fileTimestamp, err := os.Create(util.FilePathTimestamp)
	util.Check(err, "could not create timestamp file")

	_, err = fileTimestamp.Write(timestampBytes)
	util.Check(err, "could not write to timestamp file")

	err = fileTimestamp.Close()
	util.Check(err, "could not close timestamp file")
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

	err := cmd.Run()
	if err == nil {
		log.Println(outBytes.String())
	} else {
		log.Println("ERROR when running command:", err)
	}
}
