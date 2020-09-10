package util

import (
	"encoding/binary"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func Check(err error, reason string) {
	if err != nil {
		log.Fatal("ERROR ", reason, err)
	}
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

func GetTimestampBytes() []byte {
	timestampBytes := make([]byte, TimestampLen)
	binary.LittleEndian.PutUint64(timestampBytes, uint64(time.Now().UnixNano()))
	return timestampBytes
}

func GetClientSourcePort(publicKeyBytes []byte) int {
	return int(binary.LittleEndian.Uint16(publicKeyBytes))
}
