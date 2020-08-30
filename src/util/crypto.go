package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"log"
)

func GenRandomBytes(length int) [] byte {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	Check(err, "could not create random bytes")

	return bytes
}

func EncryptData(keyBytes []byte, dataBytes []byte) ([]byte, bool) {
	blockCipher, err := aes.NewCipher(keyBytes)
	if err != nil {
		log.Println("could not create cipher")
		return nil, false
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		log.Println("could not create GCM")
		return nil, false
	}

	nonceBytes := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonceBytes)
	if err != nil {
		log.Println("could not randomly generate bytes for nonce")
		return nil, false
	}

	return gcm.Seal(nonceBytes, nonceBytes, dataBytes, nil), true
}

func DecryptData(cryptoKeyBytes []byte, encryptedBytes []byte) ([]byte, bool) {
	blockCipher, err := aes.NewCipher(cryptoKeyBytes)
	if err != nil {
		log.Println("could not create cipher")
		return nil, false
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		log.Println("could not create GCM")
		return nil, false
	}

	nonceSize := gcm.NonceSize()
	if len(encryptedBytes) < nonceSize {
		log.Println("data does not match nonce size")
		return nil, false
	}

	nonceBytes, ciphertextBytes := encryptedBytes[:nonceSize], encryptedBytes[nonceSize:]
	dataBytes, err := gcm.Open(nil, nonceBytes, ciphertextBytes, nil)
	if err != nil {
		log.Println("could not decrypt")
		return nil, false
	}

	return dataBytes, true
}

