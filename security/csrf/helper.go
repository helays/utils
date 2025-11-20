package csrf

import (
	"crypto/rand"
	"encoding/base64"
	mathrand "math/rand"
	"os"
	"runtime"
	"time"
)

// GenerateCSRFToken 通用CSRF Token生成
func GenerateCSRFToken() string {
	var b [24]byte
	if _, err := rand.Read(b[:]); err != nil {
		return base64.RawURLEncoding.EncodeToString(b[:])
	}
	// 使用弱随机数生成器，但用时间戳增加熵值
	return generateMultiSourceToken()
}
func generateMultiSourceToken() string {
	data := make([]byte, 32) // 降级时用更长token
	weakRand := mathrand.New(mathrand.NewSource(time.Now().UnixNano()))

	timestamp := time.Now().UnixNano()
	pid := os.Getpid()

	for i := range data {
		// 使用可用的系统信息，避免unsafe
		data[i] = byte(weakRand.Intn(256)) ^
			byte(timestamp>>(uint(i%8)*8)) ^
			byte(pid>>(uint(i%4)*8)) ^
			byte(i) ^ // 位置信息
			byte(runtime.NumGoroutine()%256) ^ // Goroutine数量
			byte(time.Now().Nanosecond()%256) // 纳秒
	}

	return "fb_" + base64.RawURLEncoding.EncodeToString(data)
}
