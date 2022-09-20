package secret

import (
	"errors"
	"os"
)

type MockSecretProvider struct {}

func NewMockSecretProvider() *MockSecretProvider {
	return &MockSecretProvider{}
}

func (provider *MockSecretProvider) Get(name string) (string, error) {
	value := os.Getenv(name)

	if len(value) <= 0 {
		return "", errors.New("secret not found")
	}

	return value, nil
}