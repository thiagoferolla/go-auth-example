package email

import (
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendgridEmailProvider struct {
	ApiKey string
}

func NewSendgridEmailProvider(apiKey string) *SendgridEmailProvider {
	return &SendgridEmailProvider{apiKey}
}

func (provider SendgridEmailProvider) SendEmail(from string, name string, to string, templateId string, substitutions map[string]string) error {
	m := mail.NewV3Mail()

	e := mail.NewEmail(name, from)
	m.SetFrom(e)
	m.SetTemplateID(templateId)

	p := mail.NewPersonalization()

	p.AddTos(
		mail.NewEmail(name, to),
	)

	for _, key := range substitutions {
		p.SetDynamicTemplateData(key, substitutions[key])
	}

	m.AddPersonalizations(p)

	request := sendgrid.GetRequest(provider.ApiKey, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	var Body = mail.GetRequestBody(m)
	request.Body = Body
	_, err := sendgrid.API(request)

	return err
}
