package main

import (
	"bytes"
	"encoding/binary"
	"go-remote/src/client"
	"go-remote/src/server"
	"go-remote/src/util"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

func assertStartsWith(t *testing.T, actual, expected string) {
	if !strings.HasPrefix(actual, expected) {
		t.Errorf("expected actual value %v to start with %v", actual, expected)
	}
}

func assertEqual(t *testing.T, actual interface{}, expected interface{}) {
	if actual != expected {
		t.Errorf("actual value was %v, expected %v", actual, expected)
	}
}

func assertNotEqual(t *testing.T, actual interface{}, notExpected interface{}) {
	if actual == notExpected {
		t.Errorf("actual value was %v, dit not expect %v", actual, notExpected)
	}
}

func sendDataGenerator(dataToSend []byte, sourcePort int, waitTime int64) func(address string, keyFilePath string) bool {
	return func(address, keyFilePath string) bool {
		keyFileBytes := util.GetKeyBytes(keyFilePath)
		usedDataToSend := client.GetDataToSend(keyFileBytes)

		if dataToSend != nil {
			usedDataToSend = dataToSend
		}

		resolvedAddress, _ := net.ResolveUDPAddr("udp", address)

		usedSourcePort := util.GetSourcePort(keyFileBytes)

		if sourcePort > -1 {
			usedSourcePort = sourcePort
		}

		time.Sleep(time.Duration(waitTime) * time.Second)

		connection, _ := net.DialUDP("udp", &net.UDPAddr{Port: usedSourcePort}, resolvedAddress)
		defer connection.Close()
		_, _ = io.Copy(connection, bytes.NewReader(usedDataToSend))

		return false
	}
}

var currentPort = 12345

func testReceiveData(t *testing.T, keyFilePath string, timestampFileContent uint64, dataSender func(address string, keyFilePath string) bool) {
	defer os.Remove("./.timestamp")

	if keyFilePath == "" {
		keyFilePath = getKeyFilePath()
	}
	defer os.Remove(keyFilePath)

	currentPort += 1
	port := strconv.Itoa(currentPort)

	quit := make(chan bool)
	_ = os.Remove(util.FilePathTimestamp)
	go server.Run(port, keyFilePath, int64(1), "touch .start", int64(1), "touch .end", quit)

	time.Sleep(time.Millisecond)

	if timestampFileContent > 0 {
		timestampBytes := make([]byte, util.TimestampLen)
		binary.LittleEndian.PutUint64(timestampBytes, timestampFileContent)
		util.WriteBytes(util.FilePathTimestamp, timestampBytes)
	}

	success := dataSender("127.0.0.1:"+port, keyFilePath)

	quit <- true

	startFile := "./.start"
	defer os.Remove(startFile)

	endFile := "./.end"
	defer os.Remove(endFile)

	_, startErr := os.Stat(startFile)
	_, endErr := os.Stat(endFile)
	if success {
		assertEqual(t, startErr, nil)
		assertEqual(t, endErr, nil)
	} else {
		assertNotEqual(t, startErr, nil)
		assertNotEqual(t, endErr, nil)
	}
}

func getKeyFilePath() string {
	keyName := client.Run(true, "", "")
	return "./" + keyName + util.KeySuffix
}
