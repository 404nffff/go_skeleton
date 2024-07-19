package common

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5(s string) string {
	sum := md5.Sum([]byte(s))
	return hex.EncodeToString(sum[:])
}
