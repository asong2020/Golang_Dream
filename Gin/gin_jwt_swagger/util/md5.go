package util

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"

	uuid "github.com/satori/go.uuid"
)

//md5 加密算法
func MD5V(str []byte) string {
	h := md5.New()
	h.Write(str)
	return hex.EncodeToString(h.Sum(nil))
}

// uuid
func GenerateSalt() string {
	u := uuid.NewV4()
	res := fmt.Sprintf("%s", u)
	return res
}
