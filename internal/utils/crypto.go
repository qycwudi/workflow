package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
)

func Md5(str string) string {
	h := md5.New()
	_, _ = h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func Sha256(str string) string {
	h := sha256.New()
	_, _ = h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
