package utils

import (
	"crypto/md5"
	"fmt"
)

func GetMd5(str string) string {
	data := md5.Sum([]byte(str))
	return fmt.Sprintf("%x", data)
}
