package client

import (
	"bytes"
	"go-remote/src/util"
	"io"
	"log"
	"net"
	"strconv"
	"time"
)

func Run(doGenKey bool, keyfilePath string, address string) string {
	if doGenKey {
		return genKey()
	} else if keyfilePath != "" && address != "" {
		return sendData(address, keyfilePath)
	} else {
		log.Fatal("ERROR: no valid client flag combination. ",
			"Please provide either 'gen-key' to create a keypair or provide 'key-id' and 'address'")
	}

	return ""
}

func genKey() string {
	nanoSecString := strconv.FormatInt(time.Now().UnixNano(), 10)

	filePathKey := "./" + nanoSecString + util.KeySuffix

	aesKeyBytes := util.GenRandomBytes(util.AesKeySize)
	util.WriteBytes(filePathKey, aesKeyBytes)
	log.Println("Wrote key to", filePathKey)

	return nanoSecString
}

func sendData(address, keyFilePath string) string {
	keyFileBytes := util.ReadKeyBytes(keyFilePath)
	dataToSend := GetDataToSend(keyFileBytes)

	resolvedAddress, err := net.ResolveUDPAddr("udp", address)
	util.Check(err, "could not resolve address")

	connection, err := net.DialUDP("udp", &net.UDPAddr{Port: util.GetSourcePort(keyFileBytes)}, resolvedAddress)
	util.Check(err, "could not connect to udp server")

	_, err = io.Copy(connection, bytes.NewReader(dataToSend))
	util.Check(err, "could not send bytes to udp server")

	err = connection.Close()
	util.Check(err, "could not close udp connection")

	return string(dataToSend)
}

func GetDataToSend(keyFileBytes []byte) []byte {
	dataBytes := util.GetTimestampNowBytes()

	aeadKey, err := util.GetAesGcmAEAD(keyFileBytes[0:util.AesKeySize])
	util.Check(err, "")

	return util.EncryptData(aeadKey, dataBytes)
}
