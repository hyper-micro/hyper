package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrap(t *testing.T) {
	var (
		wrapErr  = errors.New("wrapErr")
		emptyErr error
		err      = errors.New("err")
	)

	assert.Equal(t, Wrap(emptyErr, err), err)
	assert.Equal(t, Wrap(wrapErr, err).Error(), "err: wrapErr")
}
