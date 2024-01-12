package errors

import "fmt"

func Wrap(wrapErr, err error) error {
	if wrapErr == nil {
		return err
	}
	return fmt.Errorf("%v: %w", err, wrapErr)
}
