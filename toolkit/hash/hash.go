package hash

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
)

func Md5(s string) string {
	m := md5.New()
	m.Write([]byte(s))
	return hex.EncodeToString(m.Sum(nil))
}

func Sha256(message string) string {
	bytes2 := sha256.Sum256([]byte(message))
	hashcode2 := hex.EncodeToString(bytes2[:])
	return hashcode2
}
