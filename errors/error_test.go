package errors

import (
	stdErrors "errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrors_New(t *testing.T) {
	err := New(101, "error")
	stdErr := stdErrors.New("error")
	assert.Equal(t, err.Error(), stdErrors.New("error").Error())
	assert.Equal(t, err.Error(), stdErr.Error())
}

func TestErrors_Is(t *testing.T) {
	err := New(100, "error")
	_, ok := Is(err)
	assert.True(t, ok)
}

func TestErrors_Error(t *testing.T) {
	err, ok := Is(New(100, "error"))
	assert.True(t, ok)
	assert.Equal(t, err.Code(), 100)
	assert.Equal(t, err.Error(), "error (code:100)")
}
