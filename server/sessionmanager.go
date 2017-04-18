package main

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var mySignUpKey = []byte("secret")

func generateToken(username string) string {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * 24)

	tokenString, err := token.SignedString(mySignUpKey)
	chk(err)
	return tokenString
}

func validateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}
