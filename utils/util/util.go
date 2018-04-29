package util

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"
)

func md5s(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// GetGUID return a 32 bytes guid
func GetGUID() string {
	b := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}

	return md5s(base64.URLEncoding.EncodeToString(b))
}

// GetRand8 return 8 byte string
func GetRand8() string {
	s32 := GetGUID()
	return s32[:8]
}
