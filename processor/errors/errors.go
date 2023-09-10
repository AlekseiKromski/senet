package errors

type ApiError struct {
	ErrorMessage string `json:"message"`
}

func NewApiErrorMessage(e error) *ApiError {
	return &ApiError{
		ErrorMessage: e.Error(),
	}
}
