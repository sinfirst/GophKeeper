package models

import "time"

const (
	Login  string = "LOGIN"
	Text   string = "TEXT"
	Card   string = "CARD"
	Binary string = "BINARY"
)

type Record struct {
	Id         int
	TypeRecord string
	Data       []byte
	Meta       string
}

type LoginJSON struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type CardJSON struct {
	Number string `json:"number"`
	Date   string `json:"date"`
	CVV    string `json:"cvv"`
}

type TokenSettings struct {
	TokenExp  time.Duration
	SecretKey string
}

type VersionBuild struct {
	Version string
	Date    string
}
