package tool

import (
	"crypto/md5"
	"encoding/hex"
)

func StrMD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

func CalculateMD5(data []byte) string {
	if data == nil {
		return ""
	}
	hash := md5.New()
	hash.Write(data)
	return hex.EncodeToString(hash.Sum(nil))
}
