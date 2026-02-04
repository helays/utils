package password

import (
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"hash"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

// PBKDF2Hasher PBKDF2 哈希器
type PBKDF2Hasher struct {
	cfg PBKDF2Config
}

func (h *PBKDF2Hasher) Hash(password string) (string, error) {
	// 生成随机盐
	salt := make([]byte, h.cfg.SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// 选择哈希函数
	var hashFunc func() hash.Hash
	switch strings.ToLower(h.cfg.HashFunc) {
	case "sha1":
		hashFunc = sha1.New
	case "sha256":
		hashFunc = sha256.New
	case "sha512":
		hashFunc = sha512.New
	default:
		hashFunc = sha256.New // 默认使用 SHA256
	}

	// 使用 PBKDF2 生成哈希
	key := pbkdf2.Key([]byte(password), salt, h.cfg.Iterations, h.cfg.KeyLength, hashFunc)

	// 编码为字符串
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(key)

	encoded := fmt.Sprintf("$pbkdf2-%s$i=%d$%s$%s",
		strings.ToLower(h.cfg.HashFunc), h.cfg.Iterations, b64Salt, b64Hash)

	return encoded, nil
}

func (h *PBKDF2Hasher) Compare(password, hashedPassword string) error {
	parts := strings.Split(hashedPassword, "$")
	if len(parts) != 5 {
		return errors.New("invalid pbkdf2 hash format")
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
	var iterations int
	var hashFunc string
	_, err = fmt.Sscanf(parts[1], "pbkdf2-%s", &hashFunc)
	if err != nil {
		return err
	}
	_, err = fmt.Sscanf(parts[2], "i=%d", &iterations)
	if err != nil {
		return err
	}

	// 选择哈希函数
	var newHashFunc func() hash.Hash
	switch strings.ToLower(hashFunc) {
	case "sha1":
		newHashFunc = sha1.New
	case "sha256":
		newHashFunc = sha256.New
	case "sha512":
		newHashFunc = sha512.New
	default:
		return fmt.Errorf("unsupported hash function: %s", hashFunc)
	}

	// 重新计算哈希
	newHash := pbkdf2.Key([]byte(password), salt, iterations, len(storedHash), newHashFunc)

	// 使用恒定时间比较
	if subtle.ConstantTimeCompare(storedHash, newHash) == 1 {
		return nil
	}
	return ErrMismatchedHashAndPassword
}
