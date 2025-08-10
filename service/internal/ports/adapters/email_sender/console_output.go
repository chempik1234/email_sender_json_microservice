package email_sender

import (
	"context"
	"email_microservice/internal/types"
	"fmt"
)

type EmailSenderRepositoryConsoleOutput struct {
}

func (e EmailSenderRepositoryConsoleOutput) SendEmail(ctx context.Context, data types.EmailData) error {
	fmt.Printf("New email:\n%v\n\n", data)
	return nil
}

func NewEmailSenderRepositoryConsoleOutput() (*EmailSenderRepositoryConsoleOutput, error) {
	return &EmailSenderRepositoryConsoleOutput{}, nil
}
