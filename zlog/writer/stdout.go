package writer

import (
	"io"
	"os"
)

func NewStdoutWriter() io.Writer {
	return os.Stdout
}
