package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"log"
)

func GetKeyBytes(keyBase64 string) []byte {
	keyBytes, base64Err := base64.StdEncoding.DecodeString(keyBase64)
	Check(base64Err, "could not decode base64 key "+keyBase64)

	fileSize := int64(len(keyBytes))
	if fileSize != AesKeySize {
		log.Fatal("ERROR: ", keyBase64, "should be exactly", AesKeySize, "bytes long, but was", fileSize)
	}
	return keyBytes
}

func GenRandomBytes(length int) []byte {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	Check(err, "could not create random bytes")

	return bytes
}

func GetHashFromBytes(dataBytes []byte) []byte {
	hash := HashFunction.New()
	hash.Write(dataBytes)
	return hash.Sum(nil)
}

func EncryptData(aead cipher.AEAD, dataBytes []byte) []byte {
	nonce := make([]byte, aead.NonceSize())
	_, err := rand.Read(nonce)
	Check(err, "could not fill nonce")

	return aead.Seal(nonce, nonce, dataBytes, nil)
}

func DecryptData(aead cipher.AEAD, encryptedBytes []byte) ([]byte, error) {
	nonceSize := aead.NonceSize()
	decryptedBytes, err := aead.Open(nil, encryptedBytes[0:nonceSize], encryptedBytes[nonceSize:], nil)
	if err != nil {
		log.Println("could not decrypt data:", err)
		return nil, err
	}

	return decryptedBytes, nil

}

func GetAesGcmAEAD(keyBytes []byte) (cipher.AEAD, error) {
	c, err := aes.NewCipher(keyBytes)
	if err != nil {
		log.Println("could not generate cipher", err)
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		log.Println("could not generate gcm", err)
		return nil, err
	}

	return gcm, nil
}
