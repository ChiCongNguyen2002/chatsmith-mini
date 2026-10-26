package failure

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"runtime"
)

type Code string

const (
	CodeInternal           Code = "internal_error"
	CodeServiceUnavailable Code = "service_unavailable"
	CodeBadRequest         Code = "bad_request"
	CodeUnauthorized       Code = "unauthorized"
	CodeNotFound           Code = "not_found"
)

var httpStatusByCode = map[Code]int{
	CodeInternal:           http.StatusInternalServerError,
	CodeServiceUnavailable: http.StatusServiceUnavailable,
	CodeBadRequest:         http.StatusBadRequest,
	CodeUnauthorized:       http.StatusUnauthorized,
	CodeNotFound:           http.StatusNotFound,
}

type AppError struct {
	Code Code
}

func (e *AppError) Error() string {
	return string(e.Code)
}

func (c Code) HTTPStatus() int {
	if status, ok := httpStatusByCode[c]; ok {
		return status
	}

	return http.StatusInternalServerError
}

func New(code Code) error {
	return withTrace(&AppError{Code: code})
}

func ErrWithTrace(err error) error {
	return withTrace(err)
}

func CodeOf(err error) Code {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code
	}

	return CodeInternal
}

func withTrace(err error) error {
	_, file, line, _ := runtime.Caller(2)

	return fmt.Errorf("%w at %s:%d", err, filepath.Base(file), line)
}
