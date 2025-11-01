package password

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/scrypt"
)

// ScryptHasher scrypt 哈希器
type ScryptHasher struct {
	cfg ScryptConfig
}

func (h *ScryptHasher) Hash(password string) (string, error) {
	// 生成随机盐
	salt := make([]byte, h.cfg.SaltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// 使用 scrypt 生成哈希
	hash, err := scrypt.Key([]byte(password), salt, h.cfg.N, h.cfg.R, h.cfg.P, h.cfg.KeyLen)
	if err != nil {
		return "", err
	}

	// 编码为字符串
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encoded := fmt.Sprintf("$scrypt$n=%d,r=%d,p=%d$%s$%s",
		h.cfg.N, h.cfg.R, h.cfg.P, b64Salt, b64Hash)

	return encoded, nil
}

func (h *ScryptHasher) Compare(password, hashedPassword string) error {
	parts := strings.Split(hashedPassword, "$")
	if len(parts) != 5 {
		return errors.New("scrypt哈希格式无效")
	}

	// 解码盐和哈希
	salt, err := base64.RawStdEncoding.DecodeString(parts[3])
	if err != nil {
		return err
	}

	storedHash, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return err
	}

	// 解析参数
	var n, r, p int
	_, err = fmt.Sscanf(parts[2], "n=%d,r=%d,p=%d", &n, &r, &p)
	if err != nil {
		return err
	}

	// 重新计算哈希
	newHash, err := scrypt.Key([]byte(password), salt, n, r, p, len(storedHash))
	if err != nil {
		return err
	}

	// 使用恒定时间比较
	if subtle.ConstantTimeCompare(storedHash, newHash) == 1 {
		return nil
	}
	return ErrMismatchedHashAndPassword
}
