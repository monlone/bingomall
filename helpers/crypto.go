package helper

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
)

// MD5 加密

const Salt string = "*$salt@*"

// MD5 签名
func MD5(text string) string {
	hash := md5.New()
	text = Salt + text + Salt
	hash.Write([]byte(text))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

// Sha 算签名
func SHA256(text string) string {
	hash := sha256.New()
	text = Salt + text + Salt
	hash.Write([]byte(text))
	return fmt.Sprintf("%x", hash.Sum(nil))
}
