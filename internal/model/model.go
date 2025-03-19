package model

import (
	"time"

	"github.com/google/uuid"

	"github.com/RomanAgaltsev/keyper/internal/database/queries"
)

type User struct {
	ID       uint64
	Login    string
	Password []byte
}

type Secret struct {
	ID        uuid.UUID
	Name      string
	Type      queries.SecretType
	metadata  []byte
	data      []byte
	comment   string
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uint64
}

type Secrets []Secret

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
