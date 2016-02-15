package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"os"
)

var encryptionKey = os.Getenv("CONFIGO_ENCRYPTION_KEY")

func encrypt(rawValue string) (string, error) {
	if len(encryptionKey) < 1 {
		return "", errors.New("CONFIGO_ENCRYPTION_KEY should be set in order to use `encrypt` function")
	}

	rawBytes := []byte(rawValue)

	if len(rawBytes)%aes.BlockSize != 0 {
		padding := aes.BlockSize - len(rawBytes)%aes.BlockSize
		padtext := bytes.Repeat([]byte{byte(0)}, padding)
		rawBytes = append(rawBytes, padtext...)
	}

	block, err := aes.NewCipher([]byte(encryptionKey))
	if err != nil {
		return "", err
	}
	ciphertext := make([]byte, aes.BlockSize+len(rawBytes))

	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], rawBytes)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decrypt(encodedValue string) (string, error) {
	if len(encryptionKey) < 1 {
		return "", errors.New("CONFIGO_ENCRYPTION_KEY should be set in order to use `decrypt` function")
	}

	block, err := aes.NewCipher([]byte(encryptionKey))
	if err != nil {
		return "", err
	}

	b, err := base64.StdEncoding.DecodeString(encodedValue)
	if err != nil {
		return "", err
	}

	if len(b) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}
	iv := b[:aes.BlockSize]
	b = b[aes.BlockSize:]

	if len(b)%aes.BlockSize != 0 {
		return "", errors.New("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(b, b)

	b = bytes.TrimRight(b, "\x00")

	return string(b), nil
}
