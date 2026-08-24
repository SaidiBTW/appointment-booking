package errors

import "fmt"

type APIError struct {
	StatusCode int    `json:"-"`
	Code       string `json:"code"`
	Message    string `json:"message"`
	Err        error  `json:"-"`
}

func New(statusCode int, code string, message string, err error) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
		Err:        err,
	}
}

func (e *APIError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}
func CustomError(statusCode int, code string, message string) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
	}
}

func NotFoundError(msg string) *APIError {
	return &APIError{
		StatusCode: 404,
		Code:       "not_found",
		Message:    msg,
	}
}
func BadRequestError(msg string, err error) *APIError {
	return &APIError{
		StatusCode: 400,
		Code:       "bad_request",
		Message:    msg,
		Err:        err,
	}
}

func InternalServerError(err error) *APIError {
	return &APIError{
		StatusCode: 500,
		Code:       "internal_server_error",
		Message:    "An internal server error occurred",
		Err:        err,
	}
}
