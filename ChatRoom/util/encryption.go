package util

/**
Author: Asong
create by on 2020/6/1
订阅号: Golang梦工厂
Features: 加密工具
*/

import (
	"crypto/md5"
	"encoding/hex"
)

//返回一个32位md5加密后的字符串
func Get32MD5Encode(password string)  string{
	h := md5.New()
	h.Write([]byte(password))
	return hex.EncodeToString(h.Sum(nil))
}

//返回一个16位md5 加密后的字符串
func Get16MD5Encode(password string) string {
	return Get16MD5Encode(password)[8:24]
}
