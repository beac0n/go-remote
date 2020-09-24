package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io/ioutil"
	"log"
	"os"
)

func ReadKeyBytes(filePath string) []byte {
	fileInfo, err := os.Stat(filePath)
	Check(err, "could not read file"+filePath)

	fileSize := fileInfo.Size()
	if fileSize != AesKeySize {
		log.Fatal("ERROR: ", filePath, "should be exactly", AesKeySize, "bytes long, but was", fileSize)
	}

	fileBytes, err := ioutil.ReadFile(filePath)
	Check(err, "could not read file bytes"+filePath)

	return fileBytes
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

func EncryptData(aead cipher.AEAD, dataBytes []byte) ([]byte, error) {
	nonce := make([]byte, aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		log.Println("could not fill nonce", err)
		return nil, err
	}

	return aead.Seal(nonce, nonce, dataBytes, nil), nil
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
