package failure

import (
	"fmt"
	"path/filepath"
	"runtime"
)

func ErrWithTrace(err error) error {
	_, file, line, _ := runtime.Caller(1)

	return fmt.Errorf("%w at %s:%d", err, filepath.Base(file), line)
}
