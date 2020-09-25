package main

import (
	"bou.ke/monkey"
	"fmt"
	"go-remote/src/client"
	"go-remote/src/util"
	"log"
	"os"
	"testing"
)

func TestInvalidClientParams(t *testing.T) {
	var loggedValue = ""
	defer monkey.Patch(log.Fatal, func(v ...interface{}) {
		loggedValue = fmt.Sprintf("%v", v)
	}).Unpatch()

	client.Run(false, "", "")
	assertEqual(t, loggedValue, "[ERROR: no valid client flag combination.  "+
		"Please provide either 'gen-key' to create a keypair or provide 'key-id' and 'address']")
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
	var loggedValue = ""
	defer monkey.Patch(log.Println, func(v ...interface{}) {
		loggedValue = fmt.Sprintf("%v", v)
	}).Unpatch()

	testReceiveData(t, sendDataGenerator(nil, -1, 2))

	assertStartsWith(t, loggedValue, "[ERROR got invalid timestamp.")
}

func TestReceiveTooLittleData(t *testing.T) {
	var loggedValue = ""
	defer monkey.Patch(log.Println, func(v ...interface{}) {
		loggedValue = fmt.Sprintf("%v", v)
	}).Unpatch()

	testReceiveData(t, sendDataGenerator([]byte{99}, -1, 0))

	assertEqual(t, loggedValue, "[ERROR received invalid bytes length. Expected 36 got 1]")
}

func TestReceiveTooLittleCloseData(t *testing.T) {
	var loggedValue = ""
	defer monkey.Patch(log.Println, func(v ...interface{}) {
		loggedValue = fmt.Sprintf("%v", v)
	}).Unpatch()

	testReceiveData(t, sendDataGenerator(make([]byte, util.EncryptedDataLen-1), -1, 0))

	assertEqual(t, loggedValue, "[ERROR received invalid bytes length. Expected 36 got 35]")
}

func TestReceiveTooMuchData(t *testing.T) {
	var loggedValue = ""
	defer monkey.Patch(log.Println, func(v ...interface{}) {
		loggedValue = fmt.Sprintf("%v", v)
	}).Unpatch()

	testReceiveData(t, sendDataGenerator(make([]byte, util.EncryptedDataLen+1), -1, 0))

	assertEqual(t, loggedValue, "[ERROR received invalid bytes length. Expected 36 got 37]")
}

func TestReceiveWrongData(t *testing.T) {
	var loggedValue = ""
	defer monkey.Patch(log.Println, func(v ...interface{}) {
		loggedValue = fmt.Sprintf("%v", v)
	}).Unpatch()

	testReceiveData(t, sendDataGenerator(make([]byte, util.EncryptedDataLen), -1, 0))

	assertEqual(t, loggedValue, "[could not decrypt data: cipher: message authentication failed]")
}

func TestReceiveDataWrongSourcePort(t *testing.T) {
	var loggedValue = ""
	defer monkey.Patch(log.Println, func(v ...interface{}) {
		loggedValue = fmt.Sprintf("%v", v)
	}).Unpatch()

	testReceiveData(t, sendDataGenerator(nil, 5555, 0))

	assertStartsWith(t, loggedValue, "[ERROR expected source port ")
}
