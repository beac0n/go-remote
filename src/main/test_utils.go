package main

import (
	"bytes"
	"go-remote/src/client"
	"go-remote/src/server"
	"go-remote/src/util"
	"io"
	"net"
	"os"
	"testing"
	"time"
)

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

func sendDataGenerator(dataToSend []byte, sourcePort int) func(address string, keyFilePath string) bool {
	return func(address, keyFilePath string) bool {
		keyFileBytes := util.ReadBytes(keyFilePath)
		usedDataToSend := client.GetDataToSend(keyFileBytes)

		if len(dataToSend) > 0 {
			usedDataToSend = dataToSend
		}

		resolvedAddress, err := net.ResolveUDPAddr("udp", address)
		util.Check(err, "could not resolve address")

		publicKeyBytes := util.GetPublicKeyBytesFromPrivateKeyBytes(keyFileBytes[util.AesKeySize:])
		usedSourcePort := util.GetSourcePort(publicKeyBytes)

		if sourcePort > -1 {
			usedSourcePort = sourcePort
		}

		connection, err := net.DialUDP("udp", &net.UDPAddr{Port: usedSourcePort}, resolvedAddress)
		util.Check(err, "could not connect to udp server")

		_, err = io.Copy(connection, bytes.NewReader(usedDataToSend))
		util.Check(err, "could not send bytes to udp server")

		err = connection.Close()
		util.Check(err, "could not close udp connection")

		return false
	}
}

func testReceiveData(t *testing.T, dataSender func(address string, keyFilePath string) bool) {
	keyName := client.Run(true, "", "")

	clientFile := "./" + keyName + "." + util.ClientSuffix
	serverFile := "./" + keyName + "." + util.ServerSuffix

	port := "12345"

	quit := make(chan bool)
	go server.Run(port, serverFile, int64(10), "touch .start", int64(1), "touch .end", quit)

	time.Sleep(time.Second)

	success := dataSender("127.0.0.1:"+port, clientFile)

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

	_ = os.Remove(clientFile)
	_ = os.Remove(serverFile)
	_ = os.Remove("./.timestamp")
	_ = os.Remove(startFile)
	_ = os.Remove(endFile)
}
