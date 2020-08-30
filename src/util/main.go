package util

import (
	"io/ioutil"
	"log"
	"os"
)

func Check(err error, reason string) {
	if err != nil {
		log.Fatal("ERROR ", reason, err)
	}
}

func GetFileBytes(fileName string, suffix string) []byte {
	keyFileBytes, err := ioutil.ReadFile("./" + fileName + "." + suffix)
	Check(err, "could not read " + suffix + " file")

	return keyFileBytes
}

func PersistData(fileName string, firstBatchBytes []byte, secondBatchBytes []byte, suffix string) string {
	filePath := "./" + fileName + "." + suffix

	file, err := os.Create(filePath)
	Check(err, "could not create "+suffix+" file")

	_, err = file.Write(append(firstBatchBytes, secondBatchBytes...))
	Check(err, "could not write to "+suffix+" file")

	err = file.Close()
	Check(err, "could not close "+suffix+" file")

	return filePath
}
