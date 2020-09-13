package client

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"go-remote/src/util"
	"io"
	"log"
	"net"
	"strconv"
	"time"
)

func Run(doGenKey *bool, keyfilePath *string, address *string) {
	if *doGenKey {
		genKeyPair()
	} else if *keyfilePath != "" && *address != "" {
		sendData(address, keyfilePath)
	} else {
		log.Fatal("no valid client flag combination. " +
			"Please provide either 'gen-key' to create a keypair or provide 'key-id' and 'address'")
	}

}

func genKeyPair() {
	nanoSecString := strconv.FormatInt(time.Now().UnixNano(), 10)

	privateKey, err := rsa.GenerateKey(rand.Reader, util.RsaKeySize)
	util.Check(err, "could not generate private key")

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	publicKeyBytes := x509.MarshalPKCS1PublicKey(&privateKey.PublicKey)

	aesKeyBytes := util.GenRandomBytes(util.AesKeySize)

	filePathPrefix := "./" + nanoSecString + "."

	filePathClientKey := filePathPrefix + util.ClientSuffix
	util.WriteBytes(filePathClientKey, append(aesKeyBytes, privateKeyBytes...))

	filePathServerKey := filePathPrefix + util.ServerSuffix
	util.WriteBytes(filePathServerKey, append(aesKeyBytes, publicKeyBytes...))

	log.Println("Wrote key pair to " + filePathServerKey + " and " + filePathClientKey)
}

func sendData(address, keyFilePath *string) {
	keyFileBytes := util.ReadBytes(*keyFilePath)
	dataToSend := getDataToSend(keyFileBytes)

	resolvedAddress, err := net.ResolveUDPAddr("udp", *address)
	util.Check(err, "could not resolve address")

	publicKeyBytes := util.GetPublicKeyBytesFromPrivateKeyBytes(keyFileBytes[util.AesKeySize:])

	connection, err := net.DialUDP("udp", &net.UDPAddr{Port: util.GetSourcePort(publicKeyBytes)}, resolvedAddress)
	util.Check(err, "could not connect to udp server")

	_, err = io.Copy(connection, bytes.NewReader(dataToSend))
	util.Check(err, "could not send bytes to udp server")

	err = connection.Close()
	util.Check(err, "could not close udp connection")
}

func getDataToSend(keyFileBytes []byte) []byte {
	dataBytes := append(util.GetTimestampNowBytes(), util.GenRandomBytes(util.SaltLen)...)

	privateKey, err := x509.ParsePKCS1PrivateKey(keyFileBytes[util.AesKeySize:])
	util.Check(err, "could not parse private key bytes")

	signedDataBytes, err := util.SignData(privateKey, dataBytes)
	util.Check(err, "")

	aeadKey, err := util.GetAesGcmAEAD(keyFileBytes[0:util.AesKeySize])
	util.Check(err, "")

	binaryHashKeyBytes := util.GetBinaryHashKeyBytes()

	aeadBinaryFirst, err := util.GetAesGcmAEAD(binaryHashKeyBytes[0:util.AesKeySize])
	util.Check(err, "")

	aeadBinarySecond, err := util.GetAesGcmAEAD(binaryHashKeyBytes[util.AesKeySize : util.AesKeySize+util.AesKeySize])
	util.Check(err, "")

	encryptedData, err := util.EncryptData(aeadKey, append(dataBytes, signedDataBytes...))
	util.Check(err, "")

	encryptedData, err = util.EncryptData(aeadBinaryFirst, encryptedData)
	util.Check(err, "")

	encryptedData, err = util.EncryptData(aeadBinarySecond, encryptedData)
	util.Check(err, "")

	return encryptedData
}
