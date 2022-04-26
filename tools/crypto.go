package tools

import (
	"crypto/sha1"
	"encoding/hex"
)

func Sha1(data string) string {
	hash := sha1.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum([]byte("")))
}
