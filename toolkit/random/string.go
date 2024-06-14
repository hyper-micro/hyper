package random

import (
	"crypto/md5"
	cRand "crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/google/uuid"
)

var (
	visibleLetters = "23456789abcdefghjkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ"
	normalLetters  = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func randString(letters string, l int) string {
	var s []string
	b := new(big.Int).SetInt64(int64(len(letters)))
	for i := 0; i < l; i++ {
		if i, err := cRand.Int(cRand.Reader, b); err == nil {
			s = append(s, string(letters[i.Int64()]))
		}
	}
	return strings.Join(s, "")
}

func String(len int) string {
	return randString(normalLetters, len)
}

func VisibleString(len int) string {
	return randString(visibleLetters, len)
}

func UUID() string {
	u, err := uuid.NewRandom()
	if err != nil {
		u, _ = uuid.NewUUID()
	}
	return u.String()
}

func Md5() string {
	m := md5.New()
	m.Write([]byte(UUID()))
	return hex.EncodeToString(m.Sum(nil))
}

func Sha256() string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(UUID())))
}
