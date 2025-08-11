package main

import (
	"context"
	"email_microservice/internal/config"
	"email_microservice/internal/ports/adapters/email_data_receiver"
	"email_microservice/internal/ports/adapters/email_sender"
	"email_microservice/internal/services"
	"email_microservice/pkg/logger"
	"email_microservice/pkg/rabbitmq"
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

	cfg, err := config.New()
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "failed to load config from env vars", zap.Error(err))
	}

	var rabbitmqManager *rabbitmq.QueueManager
	rabbitmqManager, err = rabbitmq.NewQueueManager(cfg.RabbitMQConfig)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "failed to connect to rabbitmq", zap.Error(err))
	}

	var receiverRepo *email_data_receiver.EmailDataReceiverRepositoryRabbitMQ
	receiverRepo, err = email_data_receiver.NewEmailDataReceiverRepositoryRabbitMQ(rabbitmqManager)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "failed to initialize email_data_receiver repository", zap.Error(err))
	}

	var senderRepo *email_sender.EmailSenderRepositoryNetSMTP
	senderRepo, err = email_sender.NewEmailSenderRepositoryNetSMTP(cfg.EmailConfig)
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
