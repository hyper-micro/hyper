package errors

import "fmt"

type Errors struct {
	code int
	msg  string
}

func New(code int, msg string) error {
	return &Errors{
		code: code,
		msg:  msg,
	}
}

func Is(err error) (*Errors, bool) {
	if err == nil {
		return nil, false
	}
	wrapErr, ok := err.(*Errors)
	return wrapErr, ok
}

func (e *Errors) Error() string {
	return fmt.Sprintf("%s (code:%d)", e.msg, e.code)
}

func (e *Errors) Code() int {
	return e.code
}
