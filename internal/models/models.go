package models

import "time"

type DataType string

const (
	Login  DataType = "LOGIN"
	Text   DataType = "TEXT"
	Card   DataType = "CARD"
	Binary DataType = "BINARY"
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

type TokenSettings struct {
	TokenExp  time.Duration
	SecretKey string
}

type VersionBuild struct {
	Version string
	Date    string
}
