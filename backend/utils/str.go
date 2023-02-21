// @Author HatsuneMona 2022/11/11 0:07
package utils

import (
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
