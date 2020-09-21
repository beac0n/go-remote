package main

import (
	"go-remote/src/client"
	"go-remote/src/server"
	"go-remote/src/util"
	"os"
	"testing"
	"time"
)

func Assert(t *testing.T, actual interface{}, expected interface{}) {
	if actual != expected {
		t.Errorf("actual value was %v, expected %v", actual, expected)
	}
}

func AssertComp(t *testing.T, actual interface{}, expected interface{}, comparision bool) {
	if !comparision {
		t.Errorf("actual value was %v, expected %v", actual, expected)
	}
}

func TestKeyGen(t *testing.T) {
	keyName := client.Run(true, "", "")

	serverFile := "./" + keyName + "." + util.ServerSuffix
	fileInfo, err := os.Stat(serverFile)
	Assert(t, err, nil)
	Assert(t, fileInfo.Size(), int64(util.AesKeySize+526))

	clientFile := "./" + keyName + "." + util.ClientSuffix
	fileInfo, err = os.Stat(clientFile)
	Assert(t, err, nil)
	AssertComp(t, fileInfo.Size(), ">= 2379", fileInfo.Size() >= int64(util.AesKeySize+2347))

	_ = os.Remove(serverFile)
	_ = os.Remove(clientFile)
}

func TestSendData(t *testing.T) {
	keyName := client.Run(true, "", "")
	clientFile := "./" + keyName + "." + util.ClientSuffix
	result := client.Run(false, clientFile, "127.0.0.1:8080")
	Assert(t, len(result), util.EncryptedDataLen)

	_ = os.Remove("./" + keyName + "." + util.ServerSuffix)
	_ = os.Remove(clientFile)
}

func TestReceiveData(t *testing.T) {
	keyName := client.Run(true, "", "")

	clientFile := "./" + keyName + "." + util.ClientSuffix
	serverFile := "./" + keyName + "." + util.ServerSuffix

	port := "12345"

	quit := make(chan bool)
	go server.Run(port, serverFile, int64(10), "touch .start", int64(1), "touch .end", quit)

	time.Sleep(time.Second)

	client.Run(false, clientFile, "127.0.0.1:"+port)

	quit <- true

	startFile := "./.start"
	endFile := "./.end"

	_, err := os.Stat(startFile)
	Assert(t, err, nil)

	_, err = os.Stat(endFile)
	Assert(t, err, nil)

	_ = os.Remove(clientFile)
	_ = os.Remove(serverFile)
	_ = os.Remove("./.timestamp")
	_ = os.Remove(startFile)
	_ = os.Remove(endFile)
}
