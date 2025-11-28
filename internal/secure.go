package internal

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	mathrand "math/rand"
	"strings"
	"time"
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

	EncryptedByte := []byte(base64.StdEncoding.EncodeToString(ciphertext))
	return EncryptedByte, nil
}

func Decrypt(encryption_key []byte, data string) (string, error) {

	ciphertext, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(encryption_key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// GenerateRandomKey generates a random alphanumeric key
func GenerateRandomKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	const keyLength = 16

	r := mathrand.New(mathrand.NewSource(time.Now().UnixNano()))
	var sb strings.Builder
	sb.Grow(keyLength)

	for i := 0; i < keyLength; i++ {
		sb.WriteByte(charset[r.Intn(len(charset))])
	}

	return sb.String()
}
