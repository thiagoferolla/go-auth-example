package email

type EmailProvider interface {
	SendEmail(from string, name string, to string, templateId string, substitutions map[string]string) error
}
