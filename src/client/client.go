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

func Run(doGenKey bool, keyfilePath string, address string) string {
	if doGenKey {
		return genKeyPair()
	} else if keyfilePath != "" && address != "" {
		return sendData(address, keyfilePath)
	} else {
		log.Fatal("ERROR: no valid client flag combination. ",
			"Please provide either 'gen-key' to create a keypair or provide 'key-id' and 'address'")
	}

	return ""
}

func genKeyPair() string {
	nanoSecString := strconv.FormatInt(time.Now().UnixNano(), 10)

	privateKey, err := rsa.GenerateKey(rand.Reader, util.RsaKeySize)
	util.Check(err, "could not generate private key")

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	publicKeyBytes := x509.MarshalPKCS1PublicKey(&privateKey.PublicKey)

	aesKeyBytes := util.GenRandomBytes(util.AesKeySize)

	filePathPrefix := "./" + nanoSecString

	filePathClientKey := filePathPrefix + util.ClientSuffix
	util.WriteBytes(filePathClientKey, append(aesKeyBytes, privateKeyBytes...))

	filePathServerKey := filePathPrefix + util.ServerSuffix
	util.WriteBytes(filePathServerKey, append(aesKeyBytes, publicKeyBytes...))

	log.Println("Wrote key pair to", filePathServerKey, "and", filePathClientKey)

	return nanoSecString
}

func sendData(address, keyFilePath string) string {
	keyFileBytes := util.ReadBytes(keyFilePath)
	dataToSend := GetDataToSend(keyFileBytes)

	resolvedAddress, err := net.ResolveUDPAddr("udp", address)
	util.Check(err, "could not resolve address")

	publicKeyBytes := util.GetPublicKeyBytesFromPrivateKeyBytes(keyFileBytes[util.AesKeySize:])

	connection, err := net.DialUDP("udp", &net.UDPAddr{Port: util.GetSourcePort(publicKeyBytes)}, resolvedAddress)
	util.Check(err, "could not connect to udp server")

	_, err = io.Copy(connection, bytes.NewReader(dataToSend))
	util.Check(err, "could not send bytes to udp server")

	err = connection.Close()
	util.Check(err, "could not close udp connection")

	return string(dataToSend)
}

func GetDataToSend(keyFileBytes []byte) []byte {
	dataBytes := util.GetTimestampNowBytes()

	privateKey, err := x509.ParsePKCS1PrivateKey(keyFileBytes[util.AesKeySize:])
	util.Check(err, "could not parse private key bytes")

	signedDataBytes, err := util.SignData(privateKey, dataBytes)
	util.Check(err, "")

	aeadKey, err := util.GetAesGcmAEAD(keyFileBytes[0:util.AesKeySize])
	util.Check(err, "")

	encryptedData, err := util.EncryptData(aeadKey, append(dataBytes, signedDataBytes...))
	util.Check(err, "")

	return encryptedData
}
