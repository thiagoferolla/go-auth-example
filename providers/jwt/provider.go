package jwt

import (
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/thiagoferolla/go-auth/database/models"
)

type JWTProvider interface {
	GenerateToken(user models.User) (string, error)
	ValidateToken(token string) (JwtClaims, error)
}

type JwtClaims struct {
	ID    uuid.UUID
	Email string
	Role  string
	jwt.StandardClaims
}
