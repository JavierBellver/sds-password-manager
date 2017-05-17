package main

type sessionUnit struct {
	SessionNumber string
	username      string
}

var sessionData []sessionUnit
var currentUsername string

func getSessionUsername(tokenString string) string {
	username := ""
	for i := 0; i < len(sessionData); i++ {
		var u = sessionData[i]
		if u.SessionNumber == tokenString {
			username = u.username
		}
	}
	return username
}
