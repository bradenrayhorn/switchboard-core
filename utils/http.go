package utils

import "errors"

type HttpError struct {
	Code  int
	Error error
}

func MakeHttpError(code int, message string) *HttpError {
	return &HttpError{
		Code:  code,
		Error: errors.New(message),
	}
}
