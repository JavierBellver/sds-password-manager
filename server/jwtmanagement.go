package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var mySignUpKey = generateRandomBytes(32)

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
		token, err := jwt.Parse(strings.Split(r.Header.Get("Authorization"), " ")[1], func(token *jwt.Token) (interface{}, error) {
			return []byte(mySignUpKey), nil
		})

		if err == nil && token.Valid {
			next.ServeHTTP(w, r)
		} else {
			fmt.Println("Token Error")
		}
	})
}
