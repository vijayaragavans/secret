package internal

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

func Encrypt(key []byte, plaintext string) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	if len(ciphertext) == 0 {
		return nil, fmt.Errorf("encryption produced empty ciphertext")
	}

	fmt.Println("ciphertext: ", ciphertext)

	EncryptedByte := []byte(base64.StdEncoding.EncodeToString(ciphertext))
	return EncryptedByte, nil
}
