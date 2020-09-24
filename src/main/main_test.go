package main

import (
	"go-remote/src/client"
	"go-remote/src/util"
	"os"
	"testing"
)

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
