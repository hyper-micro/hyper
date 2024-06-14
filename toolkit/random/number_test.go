package random

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInt64(t *testing.T) {
	var sets = make(map[int64]int)
	for i := 0; i < 1000; i++ {
		sets[Int64()]++
	}
	for _, count := range sets {
		assert.Equal(t, count, 1)
	}
}

func TestRange(t *testing.T) {
	assert.Equal(t, int64(0), Range(0))
	assert.Panics(t, func() {
		Range(-1)
	})

	var sets = make(map[int64]int)
	for i := 0; i < 1000; i++ {
		sets[Range(math.MaxInt32)]++

		n := Range(int64(i))
		if n > int64(i) {
			t.Errorf("n cannot be greater than %d", i)
		}
	}
	for _, count := range sets {
		assert.Equal(t, count, 1)
	}

}
