package models

type AppError string

func (e AppError) Error() string { return string(e) }

const (
	ErrUnauthenticated AppError = "unauthenticated"
	ErrConflict        AppError = "conflict"
	ErrAccessDenied    AppError = "access denied"
	ErrNotFound        AppError = "not found"
)

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
