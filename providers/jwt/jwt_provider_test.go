package jwt

import (
	"testing"
	"time"

	jwt "github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/thiagoferolla/go-auth/database/models"
)

const testSecret = "test-signing-secret"

var testUserID = uuid.MustParse("1e0f0a1c-1111-4a1b-8a2c-3d4e5f607182")

func signWithKey(t *testing.T, key []byte) string {
	t.Helper()

	claims := JwtClaims{
		ID:    testUserID,
		Email: "user@example.com",
		Role:  "user",
		StandardClaims: jwt.StandardClaims{
			Audience:  "go-auth",
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(key)

	if err != nil {
		t.Fatalf("sign token: %v", err)
	}

	return token
}

func TestValidateTokenRejectsTokenSignedWithEmptyKey(t *testing.T) {
	provider := NewBaseProvider([]byte(testSecret))

	if _, err := provider.ValidateToken(signWithKey(t, []byte(""))); err == nil {
		t.Fatal("expected a token signed with an empty key to fail validation")
	}
}

func TestValidateTokenAcceptsTokenFromGenerateToken(t *testing.T) {
	provider := NewBaseProvider([]byte(testSecret))

	token, err := provider.GenerateToken(models.User{ID: testUserID, Email: "user@example.com", Role: "user"})

	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	claims, err := provider.ValidateToken(token)

	if err != nil {
		t.Fatalf("expected the generated token to validate, got %v", err)
	}

	if claims.ID != testUserID {
		t.Fatalf("expected subject %s, got %s", testUserID, claims.ID)
	}
}
