package main

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/scrypt"
)

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
	dk, err := scrypt.Key([]byte(psw), salt, 16384, 8, 1, 32)
	chk(err)
	return base64.URLEncoding.EncodeToString(dk), base64.URLEncoding.EncodeToString(salt)
}
