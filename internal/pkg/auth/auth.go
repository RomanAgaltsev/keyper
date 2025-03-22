// Package auth provides methods for authorization and authentication.
package auth

import (
	"context"
	"errors"
	"time"

	"github.com/RomanAgaltsev/keyper/internal/model"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"golang.org/x/crypto/bcrypt"
)

type (
	UserIDName    string
	UserLoginName string
	TokenExpName  string
)

const (
	// JWTSignAlgorithm contains JWT signing algorithm.
	JWTSignAlgorithm = "HS256"

	// UserIDClaimName contains key name of user ID in a context.
	UserIDClaimName UserIDName = "uid"

	// UserLoginClaimName contains key name of user login in a context.
	UserLoginClaimName UserLoginName = "login"

	// TokenExpClaimName contains key name of token expiration in a context.
	TokenExpClaimName TokenExpName = "exp"
)

var ErrNoUserID = errors.New("no user ID provided")

// NewAuth returns new JWTAuth.
func NewAuth(secretKey string) *jwtauth.JWTAuth {
	return jwtauth.New(JWTSignAlgorithm, []byte(secretKey), nil)
}

// NewJWTToken creates new JWT token.
func NewJWTToken(ja *jwtauth.JWTAuth, user model.User, duration time.Duration) (token jwt.Token, tokenString string, err error) {
	return ja.Encode(map[string]interface{}{
		string(UserIDClaimName):    user.ID,
		string(UserLoginClaimName): user.Login,
		string(TokenExpClaimName):  time.Now().Add(duration).Unix(),
	})
}

func GetUserUID(ctx context.Context) (int32, error) {
	uidInterface := ctx.Value(UserIDClaimName)
	if uidInterface == nil {
		return 0, ErrNoUserID
	}

	uidInt32, ok := uidInterface.(int32)
	if !ok {
		return 0, ErrNoUserID
	}

	return uidInt32, nil
}

// HashPassword generates and returns hash of a given password.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compares given password and hash.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
