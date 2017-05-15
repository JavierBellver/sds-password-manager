package main

import (
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

	var unit sessionUnit
	unit.SessionNumber = tokenString
	unit.username = username
	sessionData = append(sessionData, unit)

	chk(err)
	return tokenString
}

func validateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := jwt.Parse(strings.Split(r.Header.Get("Authorization"), " ")[1], func(token *jwt.Token) (interface{}, error) {
			return []byte(mySignUpKey), nil
		})

		currentUsername = getSessionUsername(token.Raw)
		if err == nil && token.Valid {
			next.ServeHTTP(w, r)
		} else {
			response(w, false, "User token invalid")
		}
	})
}
