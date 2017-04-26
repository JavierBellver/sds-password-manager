package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	"golang.org/x/crypto/scrypt"
)

var key = generateRandomBytes(32)

func generateRandomBytes(n int) []byte {
	b := make([]byte, n)
	_, err := rand.Read(b)
	chk(err)

	return b
}

func generateRandomString(l int) string {
	b := generateRandomBytes(l)
	s := base64.URLEncoding.EncodeToString(b)
	return s
}

func hashPassword(psw string) (string, string) {
	var salt = generateRandomBytes(32)
	hashedPassword, err := scrypt.Key([]byte(psw), salt, 16384, 8, 1, 32)
	chk(err)
	return base64.URLEncoding.EncodeToString(hashedPassword), base64.URLEncoding.EncodeToString(salt)
}

func checkHashedPassword(psw string, candidatePswd string, salt string) bool {
	var result = false
	decodedSalt, err := base64.URLEncoding.DecodeString(salt)
	chk(err)
	dk, err := scrypt.Key([]byte(candidatePswd), decodedSalt, 16384, 8, 1, 32)
	chk(err)
	var candidate string = base64.URLEncoding.EncodeToString(dk)
	if candidate == psw {
		result = true
	}
	return result
}

// Encriptar string a base 64 con AES
func encrypt(key []byte, text string) string {
	plaintext := []byte(text)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return base64.URLEncoding.EncodeToString(ciphertext)
}

//Desencriptar AES de Base64 a String
func decrypt(key []byte, cryptoText string) string {
	ciphertext, err := base64.URLEncoding.DecodeString(cryptoText)
	chk(err)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%s", ciphertext)
}
