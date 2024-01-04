package writer

import (
	"io"
	"os"
)

func NewStderrWriter() io.Writer {
	return os.Stderr
}
