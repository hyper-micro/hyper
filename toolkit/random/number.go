package random

import (
	cRand "crypto/rand"
	"math"
	"math/big"
	"math/rand"
	"time"
)

var newSimpleRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func Int64() int64 {
	r, err := cRand.Int(cRand.Reader, new(big.Int).SetInt64(math.MaxInt64))
	if err != nil {
		return newSimpleRand.Int63()
	}
	return r.Int64()
}

func Range(max int64) int64 {
	if max == 0 {
		return 0
	}
	if max < 0 || max == math.MaxInt64 {
		panic("invalid argument to Range")
	}
	rag := max + 1
	r, err := cRand.Int(cRand.Reader, new(big.Int).SetInt64(rag))
	if err != nil {
		return newSimpleRand.Int63n(rag)
	}
	return r.Int64()
}
