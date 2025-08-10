package email_data_receiver

import (
	"context"
	"email_microservice/internal/types"
	"email_microservice/pkg/logger"
	"fmt"
	"strings"
)

type EmailDataReceiverRepositoryConsoleInput struct {
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func NewEmailDataReceiverRepositoryConsoleInput() (*EmailDataReceiverRepositoryConsoleInput, error) {
	return &EmailDataReceiverRepositoryConsoleInput{}, nil
}

func (e *EmailDataReceiverRepositoryConsoleInput) Run(ctx context.Context) (chan types.EmailData, chan error) {
	e.ctx, e.cancelFunc = context.WithCancel(ctx)

	emailDataChan := make(chan types.EmailData)
	errorChan := make(chan error)

	go func() {
		defer close(emailDataChan)
		defer close(errorChan)

		inputChan := make(chan string)
		go func() {
			defer close(inputChan)

			var inputString string
			for {
				_, _ = fmt.Scan(&inputString)
				if inputString == "exit" {
					break
				}
				inputChan <- inputString
			}
		}()

		currentEmail := make(map[string]string)

	out:
		for {
			select {
			case <-e.ctx.Done():
				break out
			case inputString, ok := <-inputChan:
				if !ok {
					break out
				}

				delimiterIndex := strings.Index(inputString, "=")
				if strings.TrimSpace(inputString) == "" || delimiterIndex == -1 {
					emailDataChan <- types.NewEmail(currentEmail)
					for k := range currentEmail {
						delete(currentEmail, k)
					}
				} else {
					key := strings.TrimSpace(inputString[:delimiterIndex])
					value := strings.TrimSpace(inputString[delimiterIndex+1:])
					currentEmail[key] = value
				}
			}
		}

		logger.GetLoggerFromCtx(ctx).Info(ctx, "finished console input run")
	}()

	return emailDataChan, errorChan
}

func (e *EmailDataReceiverRepositoryConsoleInput) GracefulStop() error {
	e.cancelFunc()
	return nil
}
