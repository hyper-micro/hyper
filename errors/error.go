package errors

import "fmt"

type errors struct {
	code int
	msg  string
}

func New(code int, msg string) error {
	return &errors{
		code: code,
		msg:  msg,
	}
}

func Is(err error) (*errors, bool) {
	if err == nil {
		return nil, false
	}
	wrapErr, ok := err.(*errors)
	return wrapErr, ok
}

func (e *errors) Error() string {
	return fmt.Sprintf("%s (code:%d)", e.msg, e.code)
}

func (e *errors) Code() int {
	return e.code
}
