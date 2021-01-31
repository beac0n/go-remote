package main

import (
	"bou.ke/monkey"
	"encoding/base64"
	"fmt"
	"go-remote/src/client"
	"go-remote/src/server"
	"go-remote/src/util"
	"log"
	"os"
	"strconv"
	"testing"
	"time"
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
	base64Key := getBase64Key()

	base64KeyBytes, err := base64.StdEncoding.DecodeString(base64Key)

	assertEqual(t, err, nil)
	assertEqual(t, len(base64KeyBytes), util.AesKeySize)
}

func TestSendDataFileKey(t *testing.T) {
	keyFilePath := getBase64Key()

	result := client.Run(false, keyFilePath, "127.0.0.1:8080")
	assertEqual(t, len(result), util.EncryptedDataLen)
}

func TestSendDataBase64Key(t *testing.T) {
	result := client.Run(false, getBase64Key(), "127.0.0.1:8080")
	assertEqual(t, len(result), util.EncryptedDataLen)
}

func TestInvalidTimestamp(t *testing.T) {
	defer os.Remove(util.FilePathTimestamp)

	var loggedValue = ""
	defer monkey.Patch(log.Fatal, func(v ...interface{}) {
		loggedValue = fmt.Sprintf("%v", v)
	}).Unpatch()

	util.WriteBytes(util.FilePathTimestamp, make([]byte, 9))

	quit := make(chan bool)
	go server.Run(strconv.Itoa(23456), getBase64Key(), int64(1), quit)
	quit <- true

	assertEqual(t, loggedValue, "[ERROR: /tmp/go-remote-timestamp should be exactly 8 bytes long, but was 9]")
}

func TestReceiveData(t *testing.T) {
	_ = os.MkdirAll("/tmp/go-remote", os.ModePerm)
	_, _ = os.Create("/tmp/go-remote/start")

	testReceiveData(t, "", 0, func(address string, keyFilePath string) bool {
		client.Run(false, keyFilePath, address)
		return true
	})
}

func TestReceiveDataTimestampTooNew(t *testing.T) {
	var loggedValue = ""
	defer monkey.Patch(log.Println, func(v ...interface{}) {
		loggedValue = fmt.Sprintf("%v", v)
	}).Unpatch()

	keyBase64 := getBase64Key()
	keyFileBytes := util.GetKeyBytes(keyBase64)
	aeadKey, _ := util.GetAesGcmAEAD(keyFileBytes[0:util.AesKeySize])
	encryptedData := util.EncryptData(aeadKey, util.GetTimestampNowBytes())

	testReceiveData(t, keyBase64, uint64(time.Now().UnixNano()+1), sendDataGenerator(encryptedData, -1, 0))

	assertStartsWith(t, loggedValue, "[ERROR got invalid timestamp. Expected ")
}

func TestReceiveDataTimestampInFuture(t *testing.T) {
	var loggedValue = ""
	defer monkey.Patch(log.Fatal, func(v ...interface{}) {
		loggedValue = fmt.Sprintf("%v", v)
	}).Unpatch()

	testReceiveData(t, "", 1999999999999999999, sendDataGenerator(nil, -1, 0))

	assertEqual(t, loggedValue, "[ERROR: last timestamp must be smaller than now]")
}

func TestReceiveDataTooLate(t *testing.T) {
	var loggedValue = ""
	defer monkey.Patch(log.Println, func(v ...interface{}) {
		loggedValue = fmt.Sprintf("%v", v)
	}).Unpatch()

	testReceiveData(t, "", 0, sendDataGenerator(nil, -1, 2))

	assertStartsWith(t, loggedValue, "[ERROR timestamp not within timeframe.")
}

func TestReceiveTooLittleData(t *testing.T) {
	var loggedValue = ""
	defer monkey.Patch(log.Println, func(v ...interface{}) {
		loggedValue = fmt.Sprintf("%v", v)
	}).Unpatch()

	testReceiveData(t, "", 0, sendDataGenerator([]byte{99}, -1, 0))

	assertEqual(t, loggedValue, "[ERROR received invalid bytes length. Expected 36 got 1]")
}

func TestReceiveTooLittleCloseData(t *testing.T) {
	var loggedValue = ""
	defer monkey.Patch(log.Println, func(v ...interface{}) {
		loggedValue = fmt.Sprintf("%v", v)
	}).Unpatch()

	testReceiveData(t, "", 0, sendDataGenerator(make([]byte, util.EncryptedDataLen-1), -1, 0))

	assertEqual(t, loggedValue, "[ERROR received invalid bytes length. Expected 36 got 35]")
}

func TestReceiveTooMuchData(t *testing.T) {
	var loggedValue = ""
	defer monkey.Patch(log.Println, func(v ...interface{}) {
		loggedValue = fmt.Sprintf("%v", v)
	}).Unpatch()

	testReceiveData(t, "", 0, sendDataGenerator(make([]byte, util.EncryptedDataLen+1), -1, 0))

	assertEqual(t, loggedValue, "[ERROR received invalid bytes length. Expected 36 got 37]")
}

func TestReceiveWrongData(t *testing.T) {
	var loggedValue = ""
	defer monkey.Patch(log.Println, func(v ...interface{}) {
		loggedValue = fmt.Sprintf("%v", v)
	}).Unpatch()

	testReceiveData(t, "", 0, sendDataGenerator(make([]byte, util.EncryptedDataLen), -1, 0))

	assertEqual(t, loggedValue, "[could not decrypt data: cipher: message authentication failed]")
}

func TestReceiveDataWrongSourcePort(t *testing.T) {
	var loggedValue = ""
	defer monkey.Patch(log.Println, func(v ...interface{}) {
		loggedValue = fmt.Sprintf("%v", v)
	}).Unpatch()

	testReceiveData(t, "", 0, sendDataGenerator(nil, 5555, 0))

	assertStartsWith(t, loggedValue, "[ERROR expected source port ")
}
