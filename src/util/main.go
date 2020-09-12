package util

import (
	"encoding/binary"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func Check(err error, reason string) {
	if err == nil {
		return
	}
	if reason == "" {
		os.Exit(1)
	}

	log.Fatal("ERROR ", reason, ": ", err)
}

func ReadBytes(filePath string) []byte {
	fileBytes, err := ioutil.ReadFile(filePath)
	Check(err, "could not read file "+filePath)

	return fileBytes
}

func WriteBytes(filePath string, bytes []byte) {
	file, err := os.Create(filePath)
	Check(err, "could not create file "+filePath)

	_, err = file.Write(bytes)
	Check(err, "could not write to file "+filePath)

	err = file.Close()
	Check(err, "could not close file "+filePath)
}

func GetTimestampNowBytes() []byte {
	timestampBytes := make([]byte, TimestampLen)
	binary.LittleEndian.PutUint64(timestampBytes, uint64(time.Now().UnixNano()))
	return timestampBytes
}

func GetSourcePort(publicKeyBytes []byte) int {
	sourcePort := int(binary.LittleEndian.Uint16(publicKeyBytes))
	if sourcePort < 1024 {
		return 1024
	}

	return sourcePort
}

func GetBinaryHashKeyBytesFirst() []byte {
	return getBinaryHashKeyBytes()[0:AesKeySize]
}

func GetBinaryHashKeyBytesSecond() []byte {
	return getBinaryHashKeyBytes()[AesKeySize : AesKeySize+AesKeySize]
}

func getBinaryHashKeyBytes() []byte {
	executablePath, err := os.Executable()
	Check(err, "can't get executable path of binary")

	executableBytes := ReadBytes(executablePath)
	hashBytes := GetHashFromBytes(executableBytes)
	return hashBytes
}
