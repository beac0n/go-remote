package main

import (
	"go-remote/src/client"
	"go-remote/src/util"
	"os"
	"testing"
)

func TestKeyGen(t *testing.T) {
	keyName := client.Run(true, "", "")

	serverFile := "./" + keyName + "." + util.ServerSuffix
	fileInfo, err := os.Stat(serverFile)
	assertEqual(t, err, nil)
	assertEqual(t, fileInfo.Size(), int64(util.AesKeySize+526))

	clientFile := "./" + keyName + "." + util.ClientSuffix
	fileInfo, err = os.Stat(clientFile)
	assertEqual(t, err, nil)
	if fileInfo.Size() < int64(util.AesKeySize+2347) {
		t.Errorf("actual value was %v, expected %v", fileInfo.Size(), ">= 2379")
	}

	_ = os.Remove(serverFile)
	_ = os.Remove(clientFile)
}

func TestSendData(t *testing.T) {
	keyName := client.Run(true, "", "")
	clientFile := "./" + keyName + "." + util.ClientSuffix
	result := client.Run(false, clientFile, "127.0.0.1:8080")
	assertEqual(t, len(result), util.EncryptedDataLen)

	_ = os.Remove("./" + keyName + "." + util.ServerSuffix)
	_ = os.Remove(clientFile)
}

func TestReceiveData(t *testing.T) {
	testReceiveData(t, func(address string, keyFilePath string) bool {
		client.Run(false, keyFilePath, address)
		return true
	})
}

func TestReceiveTooLittleData(t *testing.T) {
	testReceiveData(t, sendDataGenerator([]byte{99}, -1))
}

func TestReceiveTooLittleCloseData(t *testing.T) {
	testReceiveData(t, sendDataGenerator(make([]byte, util.EncryptedDataLen-1), -1))
}

func TestReceiveTooMuchData(t *testing.T) {
	testReceiveData(t, sendDataGenerator(make([]byte, util.EncryptedDataLen+1), -1))
}

func TestReceiveWrongData(t *testing.T) {
	testReceiveData(t, sendDataGenerator(make([]byte, util.EncryptedDataLen), -1))
}

func TestReceiveDataWrongSourcePort(t *testing.T) {
	testReceiveData(t, sendDataGenerator(nil, 5555))
}
