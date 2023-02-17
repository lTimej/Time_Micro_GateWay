package utils

import (
	"fmt"
	"strconv"
	"strings"
)

// 接口转string
func GetString(v interface{}) string {
	switch t := v.(type) {
	case string:
		return t
	case []string:
		return strings.Join(t, "")
	case []byte:
		return string(t)
	default:
		if v != nil {
			return fmt.Sprint(t)
		}
	}
	return ""
}

// 接口转int
func GetInt(v interface{}) int {
	switch t := v.(type) {
	case int:
		return t
	case int32:
		return int(t)
	case int64:
		return int(t)
	default:
		if d := GetString(t); d != "" {
			r, _ := strconv.Atoi(d)
			return r
		}
	}
	return 0
}
