package util

import (
	"encoding/binary"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func WriteTimestampFile(timestampBytes []byte) {
	size := len(timestampBytes)
	if size != 8 {
		log.Fatal("ERROR: timestamp should be exactly 8 bytes long, but was", size)
	}

	WriteBytes(FilePathTimestamp, timestampBytes)
}

func ReadTimestampFile() ([]byte, error) {
	fileInfo, err := os.Stat(FilePathTimestamp)
	if err != nil {
		return nil, err
	}

	fileSize := fileInfo.Size()
	if fileSize != 8 {
		log.Fatal("ERROR: ", FilePathTimestamp, " should be exactly 8 bytes long, but was", fileSize)
	}

	return ioutil.ReadFile(FilePathTimestamp)
}

func InitTimestampFile() {
	_, err := ReadTimestampFile()
	if err != nil {
		WriteTimestampFile(GetTimestampNowBytes())
	}

}

func GetTimestampNowBytes() []byte {
	timestampBytes := make([]byte, TimestampLen)
	binary.LittleEndian.PutUint64(timestampBytes, uint64(time.Now().UnixNano()))
	return timestampBytes
}
