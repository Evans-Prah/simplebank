package token

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	InvalidTokenError = fmt.Errorf("token is invalid")
	ExpiredTokenError = fmt.Errorf("token has expired")
)

type TokenPayload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewTokenPayload(username string, duration time.Duration) (*TokenPayload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &TokenPayload{
		ID:        tokenID,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, nil
}

func (payload *TokenPayload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ExpiredTokenError
	}
	return nil
}
