// @Author HatsuneMona 2022/11/11 0:07
package utils

import (
	"bytes"
	"math/rand"
	"time"
)

// RandString 生成指定长度的随机字符串
func RandString(length int) string {

	byteMap := []byte("qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM123456789")
	randLen := len(byteMap)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		bytes[i] = byteMap[r.Intn(randLen)]
	}
	return string(bytes)
}

// BuildString 通过 bytes.Buffer 创建字符串
func BuildString(s ...string) string {
	var buf bytes.Buffer

	for _, str := range s {
		buf.WriteString(str)
	}

	return buf.String()
}
