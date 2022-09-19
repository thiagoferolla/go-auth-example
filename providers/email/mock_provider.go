package email

import (
	"fmt"
	"log"
)

type MockEmailProvider struct{}

func NewMockEmailProvider() *MockEmailProvider {
	return &MockEmailProvider{}
}

func (provider *MockEmailProvider) SendEmail(from string, name string, to string, templateId string, substitutions map[string]string) error {

	log.Println(
		"Send email from ", from, " to ", fmt.Sprintf("%s | %s", name, to), " with template ", templateId,
	)

	return nil
}
