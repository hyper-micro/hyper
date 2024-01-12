package errors

import (
	stdErrors "errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrap(t *testing.T) {
	var (
		wrapErr  = stdErrors.New("wrapErr")
		emptyErr error
		err      = stdErrors.New("err")
	)

	assert.Equal(t, Wrap(emptyErr, err), err)
	assert.Equal(t, Wrap(wrapErr, err).Error(), "err: wrapErr")
}
