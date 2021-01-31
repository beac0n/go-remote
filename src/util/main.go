package util

import (
	"encoding/binary"
	"log"
	"os"
)

func Check(err error, reason string) {
	if err == nil {
		return
	} else if reason == "" {
		os.Exit(1)
	} else {
		log.Println("ERROR: ", reason, ": ", err)
		os.Exit(1)
	}
}

func GetSourcePort(keyBytes []byte) int {
	hashedPublicKeyBytes := GetHashFromBytes(keyBytes)

	sourcePort := int(binary.LittleEndian.Uint16(hashedPublicKeyBytes))
	if sourcePort < 1024 {
		return 1024
	}

	return sourcePort
}
