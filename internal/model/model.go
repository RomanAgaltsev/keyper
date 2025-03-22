package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID       int32
	Login    string
	Password string
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
	UserID    int32
}

func (s *Secret) UpdateWith(secretFrom Secret) error {
	// TODO: implement secret update logic
	// TODO: lock secret that updates
	return nil
}

type Secrets []Secret
