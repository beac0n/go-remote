package main

import (
	"bou.ke/monkey"
	"go-remote/src/client"
	"go-remote/src/util"
	"os"
	"testing"
)

func TestInvalidClientParams(t *testing.T) {
	var exitCalled = false
	defer monkey.Patch(os.Exit, func(int) {
		exitCalled = true
	}).Unpatch()

	client.Run(false, "", "")
	assertEqual(t, exitCalled, true)
}

func TestKeyGen(t *testing.T) {
	keyName := client.Run(true, "", "")

	keyFile := "./" + keyName + util.KeySuffix
	fileInfo, err := os.Stat(keyFile)
	assertEqual(t, err, nil)
	assertEqual(t, fileInfo.Size(), int64(util.AesKeySize))

	_ = os.Remove(keyFile)
}

func TestSendData(t *testing.T) {
	keyName := client.Run(true, "", "")
	keyFile := "./" + keyName + util.KeySuffix
	result := client.Run(false, keyFile, "127.0.0.1:8080")
	assertEqual(t, len(result), util.EncryptedDataLen)

	_ = os.Remove(keyFile)
}

func TestReceiveData(t *testing.T) {
	testReceiveData(t, func(address string, keyFilePath string) bool {
		client.Run(false, keyFilePath, address)
		return true
	})
}

func TestReceiveDataTooLate(t *testing.T) {
	testReceiveData(t, sendDataGenerator(nil, -1, 2))
}

func TestReceiveTooLittleData(t *testing.T) {
	testReceiveData(t, sendDataGenerator([]byte{99}, -1, 0))
}

func TestReceiveTooLittleCloseData(t *testing.T) {
	testReceiveData(t, sendDataGenerator(make([]byte, util.EncryptedDataLen-1), -1, 0))
}

func TestReceiveTooMuchData(t *testing.T) {
	testReceiveData(t, sendDataGenerator(make([]byte, util.EncryptedDataLen+1), -1, 0))
}

func TestReceiveWrongData(t *testing.T) {
	testReceiveData(t, sendDataGenerator(make([]byte, util.EncryptedDataLen), -1, 0))
}

func TestReceiveDataWrongSourcePort(t *testing.T) {
	testReceiveData(t, sendDataGenerator(nil, 5555, 0))
}
