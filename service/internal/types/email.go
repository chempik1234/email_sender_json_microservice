package types

type EmailData map[string]any

func NewEmail[T any](inputData map[string]T) EmailData {
	emailData := EmailData{}
	for k, v := range inputData {
		emailData[k] = v
	}
	return emailData
}
