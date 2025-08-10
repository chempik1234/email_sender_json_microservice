package ports

import (
	"context"
	"email_microservice/internal/types"
)

type EmailDataReceiverRepository interface {
	Run(ctx context.Context) (chan types.EmailData, chan error)
	GracefulStop() error
}

type EmailSenderRepository interface {
	SendEmail(ctx context.Context, data types.EmailData) error
}
