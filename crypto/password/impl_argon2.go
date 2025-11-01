package password

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Argon2Hasher argon2 哈希器
type Argon2Hasher struct {
	cfg  Argon2Config
	mode func(password, salt []byte, time, memory uint32, threads uint8, keyLen uint32) []byte
}

func (h *Argon2Hasher) Hash(password string) (string, error) {
	if h.mode == nil {
		h.mode = argon2.Key // 默认使用 Argon2i
	}

	// 生成随机盐
	salt := make([]byte, h.cfg.SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// 使用 Argon2 生成哈希
	hash := h.mode([]byte(password), salt, h.cfg.Iterations, h.cfg.Memory, h.cfg.Parallelism, h.cfg.KeyLength)

	// 编码为字符串: $argon2{type}$v={version}$m={memory},t={iterations},p={parallelism}${salt}${hash}
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	var algorithm HashAlgorithm
	switch h.mode {
	case argon2.IDKey:
		algorithm = HashArgon2id
	default:
		algorithm = HashArgon2
	}

	encoded := fmt.Sprintf("$%s$v=19$m=%d,t=%d,p=%d$%s$%s",
		algorithm, h.cfg.Memory, h.cfg.Iterations, h.cfg.Parallelism, b64Salt, b64Hash)

	return encoded, nil
}

func (h *Argon2Hasher) Compare(password, hashedPassword string) error {
	// 解析哈希字符串
	parts := strings.Split(hashedPassword, "$")
	if len(parts) != 6 {
		return errors.New("argon2哈希格式无效")
	}

	// 解码盐和哈希
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return err
	}

	storedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return err
	}

	// 解析参数
	var memory, iterations uint32
	var parallelism uint8
	_, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism)
	if err != nil {
		return err
	}

	// 使用相同参数重新计算哈希
	var newHash []byte
	if parts[1] == HashArgon2id.String() {
		newHash = argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, uint32(len(storedHash)))
	} else {
		newHash = argon2.Key([]byte(password), salt, iterations, memory, parallelism, uint32(len(storedHash)))
	}

	// 使用恒定时间比较
	if subtle.ConstantTimeCompare(storedHash, newHash) == 1 {
		return nil
	}
	return ErrMismatchedHashAndPassword
}
