package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	Token     uuid.UUID
	Owner     uuid.UUID
	Valid     bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewRefreshToken(owner uuid.UUID) *RefreshToken {
	return &RefreshToken{
		Token: uuid.New(),
		Owner: owner,
		Valid: true,
	}
}

type RefreshTokenRepository interface {
	GetRefreshTokenByToken(token string) (RefreshToken, error)
	InvalidateToken(token string) error
	CreateRefreshToken(refreshToken *RefreshToken, transaction *sql.Tx) (*RefreshToken, error)
}
