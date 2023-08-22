package errors

import "fmt"

type ApiError struct {
	ErrorMessage error `json:"message"`
}

func NewApiErrorMessage(e error) *ApiError {
	return &ApiError{
		ErrorMessage: e,
	}
}

func (apiError *ApiError) Error() string {
	return fmt.Sprintf(apiError.ErrorMessage.Error())
}
