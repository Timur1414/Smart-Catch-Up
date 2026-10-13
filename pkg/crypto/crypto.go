package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

var (
	ErrInvalidKeySize     = errors.New("encryption key must be exactly 32 bytes")
	ErrCiphertextTooShort = errors.New("ciphertext is too short")
)

type Encryptor struct {
	aead cipher.AEAD
}

func NewEncryptor(key []byte) (*Encryptor, error) {
	if len(key) != 32 {
		return nil, ErrInvalidKeySize
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return &Encryptor{aead: gcm}, nil
}

func (obj *Encryptor) Encrypt(plainText string) (string, error) {
	nonce := make([]byte, obj.aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	sealed := obj.aead.Seal(nonce, nonce, []byte(plainText), nil)
	return base64.StdEncoding.EncodeToString(sealed), nil
}

func (obj *Encryptor) Decrypt(cipherTextBase64 string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(cipherTextBase64)
	if err != nil {
		return "", err
	}
	nonceSize := obj.aead.NonceSize()
	if len(data) < nonceSize {
		return "", ErrCiphertextTooShort
	}
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plainBytes, err := obj.aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plainBytes), nil
}
