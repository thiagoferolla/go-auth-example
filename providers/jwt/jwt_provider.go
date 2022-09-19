package jwt

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/thiagoferolla/go-auth/database/models"
)

type JWTBaseProvider struct {
	Secret []byte
}

func NewBaseProvider() *JWTBaseProvider {
	return &JWTBaseProvider{}
}

func (provider JWTBaseProvider) GenerateToken(user models.User) (string, error) {
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

func (provider JWTBaseProvider) ValidateToken(token string) (JwtClaims, error) {
	claims := &JwtClaims{}

	t, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(provider.Secret), nil
	})

	if err != nil || !t.Valid {
		return *claims, err
	}

	return *claims, nil
}
