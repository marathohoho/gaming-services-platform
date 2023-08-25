package models

type GeneratedError struct {
	Message string `json:"message"`
}

func NewGeneratedError(message string) *GeneratedError {
	return &GeneratedError{
		Message: message,
	}
}
