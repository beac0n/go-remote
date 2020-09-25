package main

import (
	"bytes"
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
		keyFileBytes := util.ReadKeyBytes(keyFilePath)
		usedDataToSend := client.GetDataToSend(keyFileBytes)

		if dataToSend != nil {
			usedDataToSend = dataToSend
		}

		resolvedAddress, err := net.ResolveUDPAddr("udp", address)
		util.Check(err, "could not resolve address")

		usedSourcePort := util.GetSourcePort(keyFileBytes)

		if sourcePort > -1 {
			usedSourcePort = sourcePort
		}

		time.Sleep(time.Duration(waitTime) * time.Second)

		connection, err := net.DialUDP("udp", &net.UDPAddr{Port: usedSourcePort}, resolvedAddress)
		util.Check(err, "could not connect to udp server")

		_, err = io.Copy(connection, bytes.NewReader(usedDataToSend))
		util.Check(err, "could not send bytes to udp server")

		err = connection.Close()
		util.Check(err, "could not close udp connection")

		return false
	}
}

var currentPort = 12345

func testReceiveData(t *testing.T, dataSender func(address string, keyFilePath string) bool) {
	keyName := client.Run(true, "", "")

	keyFile := "./" + keyName + util.KeySuffix

	currentPort += 1
	port := strconv.Itoa(currentPort)

	quit := make(chan bool)
	go server.Run(port, keyFile, int64(1), "touch .start", int64(1), "touch .end", quit)

	time.Sleep(time.Millisecond)

	success := dataSender("127.0.0.1:"+port, keyFile)

	quit <- true

	startFile := "./.start"
	endFile := "./.end"

	_, startErr := os.Stat(startFile)
	_, endErr := os.Stat(endFile)
	if success {
		assertEqual(t, startErr, nil)
		assertEqual(t, endErr, nil)
	} else {
		assertNotEqual(t, startErr, nil)
		assertNotEqual(t, endErr, nil)
	}

	_ = os.Remove(keyFile)
	_ = os.Remove("./.timestamp")
	_ = os.Remove(startFile)
	_ = os.Remove(endFile)
}
