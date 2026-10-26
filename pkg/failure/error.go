package failure

import (
	"fmt"
	"net/http"
	"path/filepath"
	"runtime"
)

type Code string

const (
	CodeInternal           Code = "internal_error"
	CodeServiceUnavailable Code = "service_unavailable"
)

var httpStatusByCode = map[Code]int{
	CodeInternal:           http.StatusInternalServerError,
	CodeServiceUnavailable: http.StatusServiceUnavailable,
}

func (c Code) HTTPStatus() int {
	if status, ok := httpStatusByCode[c]; ok {
		return status
	}

	return http.StatusInternalServerError
}

func ErrWithTrace(err error) error {
	_, file, line, _ := runtime.Caller(1)

	return fmt.Errorf("%w at %s:%d", err, filepath.Base(file), line)
}
