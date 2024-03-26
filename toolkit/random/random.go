package random

import (
	"crypto/md5"
	cRand "crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

var newSimpleRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func UUID() string {
	u, err := uuid.NewRandom()
	if err != nil {
		u, _ = uuid.NewUUID()
	}
	return u.String()
}

func Int64() int64 {
	r, err := cRand.Int(cRand.Reader, new(big.Int).SetInt64(math.MaxInt64))
	if err != nil {
		return newSimpleRand.Int63()
	}
	return r.Int64()
}

func Md5() string {
	m := md5.New()
	m.Write([]byte(UUID()))
	return hex.EncodeToString(m.Sum(nil))
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

func Sha256() string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(UUID())))
}
