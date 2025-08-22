package models

type DataType string

const (
	Login  DataType = "LOGIN"
	Text   DataType = "TEXT"
	Card   DataType = "CARD"
	Binary DataType = "BINARY"
)

type Record struct {
	TypeRecord string
	Data       []byte
	Meta       string
}
