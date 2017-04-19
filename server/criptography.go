package main

import (
	"crypto/rand"
	"encoding/base64"
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
