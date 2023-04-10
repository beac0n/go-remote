package client

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"go-remote/src/util"
	"io"
	"log"
	"net"
)

func Run(doGenKey bool, keyBase64 string, address string) string {
	if doGenKey {
		key := genKey()
		fmt.Println("key: '" + key + "'")
		return key
	} else if keyBase64 != "" && address != "" {
		return sendData(address, keyBase64)
	} else {
		log.Fatal("ERROR: no valid client flag combination. ",
			"Please provide either 'gen-key' to create a keypair or provide 'key-id' and 'address'")
	}

	return ""
}

func genKey() string {
	aesKeyBytes := util.GenRandomBytes(util.AesKeySize)
	return base64.StdEncoding.EncodeToString(aesKeyBytes)
}

func sendData(address, keyBase64 string) string {
	keyFileBytes := util.GetKeyBytes(keyBase64)
	dataToSend := GetDataToSend(keyFileBytes)

	resolvedAddress, err := net.ResolveUDPAddr("udp", address)
	util.Check(err, "could not resolve address")

	sourcePort := util.GetSourcePort(keyFileBytes)

	if util.IsLinux() && util.IsPortInUdpSourcePorts(sourcePort) {
		log.Fatal("source port already in use:", sourcePort)
	}

	connection, err := net.DialUDP("udp", &net.UDPAddr{Port: sourcePort}, resolvedAddress)
	util.Check(err, "could not connect to udp server")

	if util.IsLinux() && !util.IsPortInUdpSourcePorts(sourcePort) {
		log.Fatal("source port was not used when creating connection:", sourcePort)
	}

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
