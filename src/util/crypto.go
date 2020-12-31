package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io/ioutil"
	"log"
	"os"
)

func GetKeyBytes(filePathOrBase64 string) []byte {
	base64KeyBytes, base64Err := getKeyBytesFromBase64(filePathOrBase64)
	fileKeyBytes, fileErr := getKeyBytesFromFile(filePathOrBase64)

	base64ErrString := "could not decode base64 key " + filePathOrBase64
	fileErrString := "could not read file " + filePathOrBase64

	if base64Err != nil && fileErr != nil {
		LogError(base64Err, base64ErrString)
		LogError(fileErr, fileErrString)
		os.Exit(1)
		return nil
	} else if base64Err != nil && fileErr == nil {
		return fileKeyBytes
	} else if base64Err == nil && fileErr != nil {
		return base64KeyBytes
	} else if base64.StdEncoding.EncodeToString(fileKeyBytes) == filePathOrBase64 {
		return base64KeyBytes
	} else {
		LogError(nil, "key is valid base64 and also a file, but they differ. Which one to choose?")
		os.Exit(1)
		return nil
	}
}

func getKeyBytesFromBase64(filePathOrBase64 string) ([]byte, error) {
	base64KeyBytes, base64Err := base64.StdEncoding.DecodeString(filePathOrBase64)
	if base64Err != nil {
		return nil, base64Err
	}

	logKeyBytesSizeError(filePathOrBase64, int64(len(base64KeyBytes)))
	return base64KeyBytes, nil
}

func getKeyBytesFromFile(filePath string) ([]byte, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	logKeyBytesSizeError(filePath, fileInfo.Size())

	fileBytes, err := ioutil.ReadFile(filePath)
	Check(err, "could not read file bytes"+filePath)

	return fileBytes, nil
}

func logKeyBytesSizeError(filePath string, fileSize int64) {
	if fileSize != AesKeySize {
		log.Fatal("ERROR: ", filePath, "should be exactly", AesKeySize, "bytes long, but was", fileSize)
	}
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
