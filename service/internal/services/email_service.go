package services

import (
	"context"
	"email_microservice/internal/ports"
	"email_microservice/internal/types"
	"email_microservice/pkg/logger"
	"go.uber.org/zap"
)

type EmailService struct {
	receiverRepo ports.EmailDataReceiverRepository
	senderRepo   ports.EmailSenderRepository
}

func NewEmailService(
	receiverRepo ports.EmailDataReceiverRepository,
	senderRepo ports.EmailSenderRepository,
) (*EmailService, error) {
	return &EmailService{
		receiverRepo: receiverRepo,
		senderRepo:   senderRepo,
	}, nil
}

func (s *EmailService) Run(ctx context.Context) error {
	emailChan, errorChan := s.receiverRepo.Run(ctx)

out:
	for {
		select {
		case err, ok := <-errorChan:
			if !ok {
				break out
			}
			if err != nil {
				logger.GetLoggerFromCtx(ctx).Error(ctx, "error while receiving emails",
					zap.Error(err))
			}
		case email, ok := <-emailChan:
			if !ok {
				break out
			}
			go s.sendEmail(ctx, email)
		case <-ctx.Done():
			_ = s.GracefulStop()
			break out
		}
	}

	logger.GetLoggerFromCtx(ctx).Info(ctx, "finished email service run")

	return s.GracefulStop()
}

func (s *EmailService) GracefulStop() error {
	return s.receiverRepo.GracefulStop()
}

func (s *EmailService) sendEmail(ctx context.Context, email types.EmailData) {
	err := s.senderRepo.SendEmail(ctx, email)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "error while sending an email",
			zap.Error(err), zap.Any("email", email))
	}
}
