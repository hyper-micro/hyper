package slice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainsString(t *testing.T) {
	assert.True(t, ContainsString("", []string{"a", "b", ""}))
	assert.False(t, ContainsString("", []string{"a", "b", "c"}))
	assert.False(t, ContainsString("d", []string{"a", "b", "c"}))
	assert.True(t, ContainsString("a", []string{"a", "b", "c"}))
}
