package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
)

func Encrypt(stringToEncrypt string, keyString string) (string, error) {
	key, err := hex.DecodeString(keyString)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	plainText := []byte(stringToEncrypt)
	cipherText := gcm.Seal(nonce, nonce, plainText, nil)

	return string(cipherText), nil

}

func Decrypt(encryptedString string, keyString string) (string, error) {
	key, err := hex.DecodeString(keyString)
	if err != nil {
		return "", err
	}

	enc, err := hex.DecodeString(encryptedString)
	if err != nil {
		return "", nil
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", nil
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", nil
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
