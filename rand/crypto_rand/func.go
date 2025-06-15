package crypto_rand

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"github.com/helays/utils/crypto/md5"
)

//
// ━━━━━━神兽出没━━━━━━
// 　　 ┏┓     ┏┓
// 　　┏┛┻━━━━━┛┻┓
// 　　┃　　　　　 ┃
// 　　┃　　━　　　┃
// 　　┃　┳┛　┗┳  ┃
// 　　┃　　　　　 ┃
// 　　┃　　┻　　　┃
// 　　┃　　　　　 ┃
// 　　┗━┓　　　┏━┛　Code is far away from bug with the animal protecting
// 　　　 ┃　　　┃    神兽保佑,代码无bug
// 　　　　┃　　　┃
// 　　　　┃　　　┗━━━┓
// 　　　　┃　　　　　　┣┓
// 　　　　┃　　　　　　┏┛
// 　　　　┗┓┓┏━┳┓┏┛
// 　　　　 ┃┫┫ ┃┫┫
// 　　　　 ┗┻┛ ┗┻┛
//
// ━━━━━━感觉萌萌哒━━━━━━
//
//
// User helay
// Date: 2024/11/12 16:27
//

// Urand 随机数生成器
func Urand() string {
	b := make([]byte, 1024)
	_, _ = rand.Read(b)
	return md5.Md5(b)
}

// SecureRandomElement 使用 crypto/rand 从任意切片或数组中安全随机返回一个元素
func SecureRandomElement[T any](collection []T) (T, bool) {
	var zero T // 用于返回类型的零值

	if len(collection) == 0 {
		return zero, false
	}

	// 生成随机索引
	randomIndex, err := secureRandomIntn(len(collection))
	if err != nil {
		// 如果 crypto/rand 失败，可以回退到 math/rand
		// 但要注意这会降低安全性
		return zero, false
	}

	return collection[randomIndex], true
}

// secureRandomIntn 使用 crypto/rand 生成 [0, n) 范围内的随机数
func secureRandomIntn(n int) (int, error) {
	if n <= 0 {
		return 0, fmt.Errorf("invalid argument to secureRandomIntn: n must be positive")
	}

	// 计算需要的字节数
	// 对于 n <= 256 只需要 1 字节，但为了简化我们总是用 4 字节
	var randomUint32 uint32
	err := binary.Read(rand.Reader, binary.BigEndian, &randomUint32)
	if err != nil {
		return 0, err
	}

	// 将随机数映射到 [0, n) 范围
	// 这种方法会有轻微偏差，但对于大多数应用可以接受
	// 如果需要完全无偏的分布，需要更复杂的算法
	result := int(randomUint32 % uint32(n))
	return result, nil
}
