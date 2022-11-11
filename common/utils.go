package common

import (
	"crypto/md5"
	"fmt"
)

// GetMd5 获取 MD5
func GetMd5(data []byte) string {
	return fmt.Sprintf("%x", md5.Sum(data))
}
