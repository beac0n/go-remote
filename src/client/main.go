package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"flag"
	"go-remote/src/util"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

func main() {
	keyFilePath := flag.String("key", "", "path to key file")
	doGenKey := flag.Bool("gen-key", false, "generate key pair")
	address := flag.String("address", "", "udp address")

	flag.Parse()

	run(doGenKey, keyFilePath, address)
}

func run(doGenKey *bool, keyfilePath *string, address *string) {
	if *doGenKey {
		genKeyPair()
	} else if *keyfilePath != "" && *address != "" {
		dataToSend, publicKeyBytes := getDataToSend(keyfilePath)
		sendData(address, dataToSend, publicKeyBytes)
	} else {
		log.Fatal("no valid client flag combination. " +
			"Please provide either 'gen-key' to create a keypair or provide 'key-id', 'address'")
	}

}

func genKeyPair() {
	nanoSecString := strconv.FormatInt(time.Now().UnixNano(), 10)

	privateKey, err := rsa.GenerateKey(rand.Reader, util.KeySize)
	util.Check(err, "could not generate private key")

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	publicKeyBytes := x509.MarshalPKCS1PublicKey(&privateKey.PublicKey)

	filePathPrefix := "./" + nanoSecString + "."

	filePathClientKey := filePathPrefix + util.ClientSuffix
	util.WriteBytes(filePathClientKey, publicKeyBytes)

	filePathServerKey := filePathPrefix + util.ServerSuffix
	util.WriteBytes(filePathServerKey, privateKeyBytes)

	log.Println("Written key pair to " + filePathServerKey + " and " + filePathClientKey)
}

func getDataToSend(keyFilePath *string) ([]byte, []byte) {
	publicKeyBytes := util.ReadBytes(*keyFilePath)

	dataBytes := append(util.GetTimestampBytes(), util.GenRandomBytes(util.SaltLen)...)
	encryptedData, success := util.EncryptData(publicKeyBytes, dataBytes)
	if !success {
		os.Exit(1)
	}

	return encryptedData, publicKeyBytes
}

func sendData(address *string, dataToSend []byte, publicKeyBytes []byte) {
	resolvedAddress, err := net.ResolveUDPAddr("udp", *address)
	util.Check(err, "could not resolve address")

	connection, err := net.DialUDP("udp", &net.UDPAddr{Port: util.GetClientSourcePort(publicKeyBytes)}, resolvedAddress)
	util.Check(err, "could not connect to udp server")

	_, err = io.Copy(connection, bytes.NewReader(dataToSend))
	util.Check(err, "could not send bytes to udp server")

	err = connection.Close()
	util.Check(err, "could not close udp connection")
}
