package util

import (
	"crypto/aes"
	"crypto/cipher"
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

func VerifySignedData(publicKey *rsa.PublicKey, dataBytes []byte, signatureBytes []byte) bool {
	hash := HashFunction.New()
	hash.Write(dataBytes)

	options := rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthAuto}
	if err := rsa.VerifyPSS(publicKey, HashFunction, hash.Sum(nil), signatureBytes, &options); err != nil {
		log.Println("could not verify signature", err)
		return false
	}

	return true
}

func SignData(privateKey *rsa.PrivateKey, dataBytes []byte) ([]byte, error) {
	hash := HashFunction.New()
	hash.Write(dataBytes)

	options := rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthAuto}
	signature, err := rsa.SignPSS(rand.Reader, privateKey, HashFunction, hash.Sum(nil), &options)
	if err != nil {
		log.Println("could not sign data", err)
		return nil, err
	}

	return signature, nil
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

func GetPublicKeyBytesFromPrivateKeyBytes(privateKeyBytes []byte) []byte {
	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBytes)
	Check(err, "could not parse private key bytes")

	return x509.MarshalPKCS1PublicKey(&privateKey.PublicKey)
}
