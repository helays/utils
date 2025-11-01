package password

import (
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// 全局默认密码管理器
var defaultPassword *Password

func init() {
	defaultPassword = NewPassword()
}

// Hash 使用默认配置哈希密码
func Hash(password string) (string, error) {
	return defaultPassword.Hash(password)
}

// Compare 使用默认配置比较密码
func Compare(password, hashedPassword string) error {
	return defaultPassword.Compare(password, hashedPassword)
}

type Password struct {
	config     *Security
	hashers    map[HashAlgorithm]Hasher
	autoDetect bool // 是否自动检测算法
}

func NewPassword(cogs ...*Security) *Password {
	cfg := defaultSecurityConfig()
	if len(cogs) > 0 && cogs[0] != nil {
		cfg = cogs[0]
	}

	cfg.ensureDefaults()

	pm := &Password{
		config:  cfg,
		hashers: make(map[HashAlgorithm]Hasher),
	}

	// 初始化各算法哈希器
	pm.hashers[HashBcrypt] = &BcryptHasher{cost: cfg.Bcrypt.Cost}
	pm.hashers[HashArgon2] = &Argon2Hasher{cfg: cfg.Argon2}
	pm.hashers[HashArgon2id] = &Argon2Hasher{cfg: cfg.Argon2id, mode: argon2.IDKey}
	pm.hashers[HashScrypt] = &ScryptHasher{cfg: cfg.Scrypt}
	pm.hashers[HashPBKDF2] = &PBKDF2Hasher{cfg: cfg.PBKDF2}

	return pm
}

// SetAutoDetect 设置是否自动检测算法
func (pm *Password) SetAutoDetect(autoDetect bool) {
	pm.autoDetect = autoDetect
}

// Hash 使用配置的算法对密码进行哈希
func (pm *Password) Hash(password string) (string, error) {
	hasher, exists := pm.hashers[pm.config.PasswordAlgorithm]
	if !exists {
		return "", fmt.Errorf("%w: %s", ErrUnsupportedAlgorithm, pm.config.PasswordAlgorithm)
	}
	return hasher.Hash(password)
}

// Compare 比较密码和哈希值
func (pm *Password) Compare(password, hashedPassword string) error {
	var (
		algorithm HashAlgorithm
		err       error
	)
	if pm.autoDetect {
		// 从哈希值中检测算法
		algorithm, err = detectAlgorithm(hashedPassword)
	} else {
		algorithm = pm.config.PasswordAlgorithm
	}
	if err != nil {
		return err
	}

	hasher, exists := pm.hashers[algorithm]
	if !exists {
		return fmt.Errorf("%w: %s", ErrUnsupportedAlgorithm, algorithm)
	}

	return hasher.Compare(password, hashedPassword)
}

// detectAlgorithm 从哈希字符串中检测算法
func detectAlgorithm(hashedPassword string) (HashAlgorithm, error) {
	if len(hashedPassword) == 0 {
		return "", errors.New("哈希字符串为空")
	}

	// 简单的前缀检测就足够了，具体格式验证交给各算法的Compare方法
	switch {
	case strings.HasPrefix(hashedPassword, "$2a$"),
		strings.HasPrefix(hashedPassword, "$2b$"),
		strings.HasPrefix(hashedPassword, "$2x$"),
		strings.HasPrefix(hashedPassword, "$2y$"):
		return HashBcrypt, nil

	case strings.HasPrefix(hashedPassword, "$argon2id$"):
		return HashArgon2id, nil

	case strings.HasPrefix(hashedPassword, "$argon2i$"),
		strings.HasPrefix(hashedPassword, "$argon2$"):
		return HashArgon2, nil

	case strings.HasPrefix(hashedPassword, "$scrypt$"):
		return HashScrypt, nil

	case strings.HasPrefix(hashedPassword, "$pbkdf2-"):
		return HashPBKDF2, nil

	default:
		return "", errors.New("无法检测哈希算法")
	}
}
