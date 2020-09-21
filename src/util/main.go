package util

import (
	"encoding/binary"
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

	log.Fatal("ERROR: ", reason, ": ", err)
}

func WriteTimestampFile(timestampBytes []byte) {
	WriteBytes(FilePathTimestamp, timestampBytes)
}

func ReadTimestampFile() ([]byte, error) {
	fileInfo, err := os.Stat(FilePathTimestamp)
	if err != nil {
		return nil, err
	}

	fileSize := fileInfo.Size()
	if fileSize != 8 {
		log.Fatal("ERROR: ", FilePathTimestamp, " should be exactly 8 bytes long, but was ", fileSize)
	}

	return ioutil.ReadFile(FilePathTimestamp)
}

func ReadBytes(filePath string) []byte {
	fileInfo, err := os.Stat(filePath)
	Check(err, "could not read file"+filePath)

	fileSize := float64(fileInfo.Size()) / 1000 / 1000
	if fileSize > MaxFileSizeMb {
		log.Fatal("expected file size to not be bigger than " +
			fmt.Sprintf("%f", MaxFileSizeMb) + " MB but was " +
			fmt.Sprintf("%f", fileSize) + " MB")
	}

	fileBytes, err := ioutil.ReadFile(filePath)
	Check(err, "could not read file bytes"+filePath)

	return fileBytes
}

func WriteBytes(filePath string, bytes []byte) {
	if fileSize := float64(len(bytes) / 1000 / 1000); fileSize > MaxFileSizeMb {
		log.Fatal("expected file size to not be bigger than " +
			fmt.Sprintf("%f", MaxFileSizeMb) + " MB but was " +
			fmt.Sprintf("%f", fileSize) + " MB")
	}

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
