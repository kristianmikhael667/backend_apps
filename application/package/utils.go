package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"os"
)

func Getenv(key, fallback string) string {
	var (
		val     string
		isExist bool
	)
	val, isExist = os.LookupEnv(key)
	if !isExist {
		val = fallback
	}
	return val
}

func GenerateBase62EncodedRandomBytes(length int) (string, error) {
	randomBytes, err := GenerateRandomBytes(length)
	if err != nil {
		return "", err
	}

	// Encode menggunakan base62
	encoded := base64.StdEncoding.EncodeToString(randomBytes)
	return encoded, nil
}

func GenerateRandomBytes(length int) ([]byte, error) {
	randomBytes := make([]byte, length)
	if _, err := rand.Read(randomBytes); err != nil {
		return nil, err
	}
	return randomBytes, nil
}

func GenerateKey(secretKey string) []byte {
	hasher := sha256.New()
	hasher.Write([]byte(Getenv("SECRET_ECRYPT", "")))
	return hasher.Sum(nil)
}

func Encrypt(plainText string, secretKey string) (string, error) {
	key := GenerateKey(secretKey)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	plaintext := []byte(plainText)
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(cipherText string, secretKey string) (string, error) {
	key := GenerateKey(secretKey)

	ciphertext, err := base64.URLEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext is too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}
