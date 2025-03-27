package client

import (
	"context"

	"github.com/google/uuid"

	"github.com/RomanAgaltsev/keyper/internal/model"
)

func NewUserClient() *UserClient {
	return nil
}

type UserClient struct {
}

func (c *UserClient) Register(ctx context.Context, user *model.User) error {
	return nil
}

func (c *UserClient) Login(ctx context.Context, user *model.User) error {
	return nil
}

func NewSecretClient() *SecretClient {
	return nil
}

type SecretClient struct {
}

func (c *SecretClient) Create(ctx context.Context, secret *model.Secret) (uuid.UUID, error) {
	return uuid.Nil, nil
}

func (c *SecretClient) Update(ctx context.Context, userID uuid.UUID, secret *model.Secret) error {
	return nil
}
func (c *SecretClient) UpdateData(ctx context.Context, userID uuid.UUID, secret *model.Secret) error {
	return nil
}

func (c *SecretClient) Get(ctx context.Context, userID uuid.UUID, secretID uuid.UUID) (*model.Secret, error) {
	return nil, nil
}

func (c *SecretClient) GetData(ctx context.Context, userID uuid.UUID, secretID uuid.UUID) (*model.Secret, error) {
	return nil, nil
}

func (c *SecretClient) List(ctx context.Context, userID uuid.UUID) (model.Secrets, error) {
	return nil, nil
}

func (c *SecretClient) Delete(ctx context.Context, userID uuid.UUID, secretID uuid.UUID) error {
	return nil
}
