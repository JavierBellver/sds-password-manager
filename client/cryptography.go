package main

import "crypto/sha256"

func hashPassword(psw string) string {
	h := sha256.New()
	h.Write([]byte(psw))
	hash := h.Sum(nil)
	return string(hash)
}
