package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	Login     string
	Password  string
	CreatedAt time.Time
}

type SecretType string

const (
	SecretTypeUNSPECIFIED SecretType = "UNSPECIFIED"
	SecretTypeCREDENTIALS SecretType = "CREDENTIALS"
	SecretTypeTEXT        SecretType = "TEXT"
	SecretTypeBINARY      SecretType = "BINARY"
	SecretTypeCARD        SecretType = "CARD"
)

type Secret struct {
	ID        uuid.UUID
	Name      string
	Type      SecretType
	Metadata  []byte
	Data      []byte
	Comment   string
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
}

func (s *Secret) UpdateWith(secretFrom *Secret) error {
	// TODO: implement secret update logic
	// TODO: lock secret that updates
	return nil
}

type Secrets []*Secret
