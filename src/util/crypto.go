package util

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"log"
)

func GenRandomBytes(length int) []byte {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	Check(err, "could not create random bytes")

	return bytes
}

func EncryptData(publicKeyBytes []byte, dataBytes []byte) ([]byte, bool) {
	publicKey, err := x509.ParsePKCS1PublicKey(publicKeyBytes)
	if err != nil {
		log.Println("could not parse public key bytes", err)
		return nil, false
	}

	encryptedBytes, err := rsa.EncryptOAEP(HashFunction.New(), rand.Reader, publicKey, dataBytes, []byte(""))
	if err != nil {
		log.Println("could not encrypt data bytes", err)
		return nil, false
	}

	return encryptedBytes, true
}

func DecryptData(privateKeyBytes []byte, encryptedBytes []byte) ([]byte, bool) {
	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBytes)
	if err != nil {
		log.Println("could not parse private key bytes", err)
		return nil, false
	}

	dataBytes, err := rsa.DecryptOAEP(HashFunction.New(), rand.Reader, privateKey, encryptedBytes, []byte(""))
	if err != nil {
		log.Println("could not decrypt encrypted bytes", err)
		return nil, false
	}

	return dataBytes, true
}
