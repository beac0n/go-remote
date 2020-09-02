package client

import (
	"bytes"
	"crypto/ed25519"
	"encoding/binary"
	"go-remote/src/util"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

func Run(doGenKey *bool, keyId *string, address *string) {
	if *doGenKey {
		genKeyPair()
	} else if *keyId != "" && *address != "" {
		dataToSend := getDataToSend(keyId)
		sendData(address, dataToSend)
	} else {
		log.Fatal("no valid client flag combination. " +
			"Please provide either 'gen-key' to create a keypair or provide 'key-id', 'address'")
	}

}

func genKeyPair() {
	serverKey, clientKey, err := ed25519.GenerateKey(nil)
	util.Check(err, "could not generate key pair")

	nanoSecString := strconv.FormatInt(time.Now().UnixNano(), 10)

	cryptoKey := util.GenRandomBytes(util.CryptoKeyLen)

	filePathServerKey := util.GetServerKeyFilePath(nanoSecString)
	util.WriteBytes(filePathServerKey, append(serverKey, cryptoKey...))

	filePathClientKey := util.GetClientKeyFilePath(nanoSecString)
	util.WriteBytes(filePathClientKey, append(clientKey, cryptoKey...))

	log.Println("Written key pair to " + filePathServerKey + " and " + filePathClientKey)
}

func getDataToSend(keyId *string) []byte {
	keyFileBytes := util.ReadBytes(util.GetClientKeyFilePath(*keyId))

	timestampBytes := make([]byte, util.TimestampLen)
	binary.LittleEndian.PutUint64(timestampBytes, uint64(time.Now().UnixNano()))

	saltBytes := util.GenRandomBytes(util.SaltLen)

	timestampBytesAndSaltBytes := append(timestampBytes, saltBytes...)
	clientKeyBytes := keyFileBytes[0:util.ClientKeyLen]

	signatureBytes := ed25519.Sign(clientKeyBytes, timestampBytesAndSaltBytes)
	dataBytes := append(timestampBytesAndSaltBytes, signatureBytes...)

	cryptoKeyBytes := keyFileBytes[util.ClientKeyLen:util.ClientKeyFileLen]
	encryptedData, success := util.EncryptData(cryptoKeyBytes, dataBytes)
	if !success {
		os.Exit(1)
	}

	return encryptedData
}

func sendData(address *string, dataToSend []byte) {
	resolvedAddress, err := net.ResolveUDPAddr("udp", *address)
	util.Check(err, "could not resolve address")

	connection, err := net.DialUDP("udp", nil, resolvedAddress)
	util.Check(err, "could not connect to udp server")

	_, err = io.Copy(connection, bytes.NewReader(dataToSend))
	util.Check(err, "could not send bytes to udp server")

	err = connection.Close()
	util.Check(err, "could not close udp connection")
}
