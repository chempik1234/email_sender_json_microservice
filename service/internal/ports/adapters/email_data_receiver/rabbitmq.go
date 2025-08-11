package email_data_receiver

import (
	"context"
	"email_microservice/internal/types"
	"email_microservice/pkg/logger"
	"email_microservice/pkg/rabbitmq"
	"encoding/json"
	"fmt"
)

type EmailDataReceiverRepositoryRabbitMQ struct {
	ctx        context.Context
	cancelFunc context.CancelFunc

	queueManager *rabbitmq.QueueManager
}

func NewEmailDataReceiverRepositoryRabbitMQ(queueManager *rabbitmq.QueueManager) (*EmailDataReceiverRepositoryRabbitMQ, error) {
	return &EmailDataReceiverRepositoryRabbitMQ{
		queueManager: queueManager,
	}, nil
}

func (e *EmailDataReceiverRepositoryRabbitMQ) Run(ctx context.Context) (chan types.EmailData, chan error) {
	e.ctx, e.cancelFunc = context.WithCancel(ctx)

	emailDataChan := make(chan types.EmailData)
	errorChan := make(chan error)

	go func() {
		defer close(emailDataChan)
		defer close(errorChan)

		inputChan, err := e.queueManager.Consume()
		if err != nil {
			errorChan <- err
			return
		}

	out:
		for {
			select {
			case <-e.ctx.Done():
				break out
			case inputDelivery, ok := <-inputChan:
				if !ok {
					break out
				}
				go func() {
					var currentEmail = make(map[string]any)
					err = json.Unmarshal(inputDelivery.Body, &currentEmail)
					if err != nil {
						errorChan <- fmt.Errorf("error while unmarshaling input body: %w", err)
					} else if len(currentEmail) > 0 {
						emailDataChan <- types.NewEmail(currentEmail)
					} else {
						errorChan <- fmt.Errorf("no email data received")
					}
				}()
			}
		}

		logger.GetLoggerFromCtx(ctx).Info(ctx, "finished rabbitmq input run")
	}()

	return emailDataChan, errorChan
}

func (e *EmailDataReceiverRepositoryRabbitMQ) GracefulStop() error {
	e.cancelFunc()
	return e.queueManager.Close()
}
