package email_sender

import (
	"context"
	"email_microservice/internal/config"
	"email_microservice/internal/types"
	"fmt"
	"net/smtp"
)

type EmailSenderRepositoryNetSMTP struct {
	config config.EmailConfig
}

func (e EmailSenderRepositoryNetSMTP) SendEmail(ctx context.Context, data types.EmailData) error {
	to := []string{e.config.To}
	from := e.config.From

	msg := "New form data!\n"
	for k, v := range data {
		msg += fmt.Sprintf("\n%s: %s\r\n", k, v)
	}

	return smtp.SendMail(
		fmt.Sprintf("%s:%d", e.config.Host, e.config.Port),
		smtp.PlainAuth("", from, e.config.Password, e.config.Host),
		from, to,
		[]byte(msg),
	)
}

func NewEmailSenderRepositoryNetSMTP(config config.EmailConfig) (*EmailSenderRepositoryNetSMTP, error) {
	return &EmailSenderRepositoryNetSMTP{config}, nil
}
