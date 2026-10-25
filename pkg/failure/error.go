package failure

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrBadRequest    = errors.New("bad request")
	ErrInvalidCursor = errors.New("invalid cursor")
)

func ErrWithTrace(err error) error {
	_, file, line, _ := runtime.Caller(1)

	return fmt.Errorf("%w at %s:%d", err, filepath.Base(file), line)
}
