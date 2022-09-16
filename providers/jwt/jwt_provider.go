package jwt

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/thiagoferolla/go-auth/database/models"
)

type JWTProvider struct {
	Secret []byte
}

func NewProvider() *JWTProvider {
	return &JWTProvider{}
}

type JwtClaims struct {
	ID    uuid.UUID
	Email string
	Role  string
	jwt.StandardClaims
}

func (provider JWTProvider) GenerateToken(user models.User) (string, error) {
	expiration := time.Now().Add(time.Second * 3700)

	claims := JwtClaims{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
		StandardClaims: jwt.StandardClaims{
			Audience:  "go-auth",
			ExpiresAt: expiration.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(provider.Secret)

	return signedToken, err
}

func (provider JWTProvider) ValidateToken(token string) (JwtClaims, error) {
	claims := &JwtClaims{}

	t, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(provider.Secret), nil
	})

	if err != nil || !t.Valid {
		return *claims, err
	}

	return *claims, nil
}
