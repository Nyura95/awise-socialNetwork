package helpers

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
)

// StringCheckFormat standardize the check string
func StringCheckFormat(str string) string {
	return strings.ToLower(strings.Trim(str, " "))
}

// StringToMD5 return str MD5
func StringToMD5(str string) string {
	hasher := md5.New()
	hasher.Write([]byte(str))
	return hex.EncodeToString(hasher.Sum(nil))
}

// Stringify return a string from interfece
func Stringify(data interface{}) string {
	return fmt.Sprintf("%b", data)
}
