package main

type sessionUnit struct {
	SessionNumber string
	username      string
}

var sessionData []sessionUnit
var currentUsername string
