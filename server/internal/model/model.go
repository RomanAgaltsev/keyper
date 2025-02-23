package model

import "time"

type User struct {
	ID       uint
	Login    string
	Password []byte
}

const (
	Account uint8 = 1 + iota
	Email
	Card
	Binary
)

type RecordType uint8

type Record struct {
	ID          uint
	Type        RecordType
	Address     string
	Credentials []Credential
	Comment     string
}

type Credential struct {
	Login     string
	Password  string
	CreatedAt time.Time
	Comment   string
}
