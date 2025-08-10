package main

import (
	"context"
	"email_microservice/internal/config"
	"email_microservice/internal/ports/adapters/email_data_receiver"
	"email_microservice/internal/ports/adapters/email_sender"
	"email_microservice/internal/services"
	"email_microservice/pkg/logger"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
)

func main() {
	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()
	ctx, _ = logger.New(ctx)

	_, err := config.New()
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "failed to load config from env vars", zap.Error(err))
	}

	var receiverRepo *email_data_receiver.EmailDataReceiverRepositoryConsoleInput
	receiverRepo, err = email_data_receiver.NewEmailDataReceiverRepositoryConsoleInput()
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "failed to initialize email_data_receiver repository", zap.Error(err))
	}

	var senderRepo *email_sender.EmailSenderRepositoryConsoleOutput
	senderRepo, err = email_sender.NewEmailSenderRepositoryConsoleOutput()
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "failed to initialize email_sender repository", zap.Error(err))
	}

	var service *services.EmailService
	service, err = services.NewEmailService(receiverRepo, senderRepo)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "failed to initialize service", zap.Error(err))
	}

	go func() {
		log.Fatal(service.Run(ctx))
	}()

	<-ctx.Done()

	_ = service.GracefulStop()
}
