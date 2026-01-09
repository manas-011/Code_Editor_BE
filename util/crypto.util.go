package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"os"
)

func Encrypt(text string) (string, error) {
	key := []byte(os.Getenv("COOKIE_SECRET")) // 32 bytes
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, _ := cipher.NewGCM(block)
	nonce := make([]byte, gcm.NonceSize())
	io.ReadFull(rand.Reader, nonce)

	ciphertext := gcm.Seal(nonce, nonce, []byte(text), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(enc string) (string, error) {
	key := []byte(os.Getenv("COOKIE_SECRET"))
	data, _ := base64.StdEncoding.DecodeString(enc)

	block, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(block)

	nonceSize := gcm.NonceSize()
	nonce, cipherText := data[:nonceSize], data[nonceSize:]

	plain, err := gcm.Open(nil, nonce, cipherText, nil)
	return string(plain), err
}
