package util

import (
	"encoding/binary"
	"errors"
	"fmt"
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

func ReadFileWithMaxSize(filePath string) ([]byte, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	fileSize := float64(fileInfo.Size()) / 1000 / 1000
	if fileSize > MaxFileSizeMb {
		return nil, errors.New("expected file size to not be bigger than " +
			fmt.Sprintf("%f", MaxFileSizeMb) + " MB but was " +
			fmt.Sprintf("%f", fileSize) + " MB")
	}

	return ioutil.ReadFile(filePath)
}

func ReadBytes(filePath string) []byte {
	fileInfo, err := os.Stat(filePath)
	Check(err, "could get file stat "+filePath)

	maxFileSizeMb := float64(3)
	fileSize := float64(fileInfo.Size()) / 1000 / 1000
	if fileSize > maxFileSizeMb {
		log.Fatal("expected file size to not be bigger than ", maxFileSizeMb, " MB but was ", fileSize, " MB")
	}

	fileBytes, err := ReadFileWithMaxSize(filePath)
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

func GetBinaryHashKeyBytes() []byte {
	executablePath, err := os.Executable()
	Check(err, "can't get executable path of binary")

	executableBytes := ReadBytes(executablePath)
	hashBytes := GetHashFromBytes(executableBytes)
	return hashBytes
}
