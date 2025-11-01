package password

import "errors"

// HashAlgorithm 密码哈希算法
type HashAlgorithm string

func (h HashAlgorithm) String() string {
	return string(h)
}

const (
	HashBcrypt   HashAlgorithm = "bcrypt"   // bcrypt 算法，安全性高，抗GPU破解，默认算法
	HashArgon2   HashAlgorithm = "argon2"   // Argon2 算法，密码哈希竞赛冠军，抗侧信道攻击
	HashArgon2id HashAlgorithm = "argon2id" // Argon2id 算法，混合模式，平衡抗GPU和侧信道攻击
	HashScrypt   HashAlgorithm = "scrypt"   // scrypt 算法，内存困难型，抗ASIC和GPU攻击
	HashPBKDF2   HashAlgorithm = "pbkdf2"   // PBKDF2 算法，经典可靠，配置灵活
)

// Security 安全配置
type Security struct {
	PasswordAlgorithm HashAlgorithm `json:"password_algorithm" yaml:"password_algorithm" ini:"password_algorithm"` // 密码hash算法
	Bcrypt            BcryptConfig  `json:"bcrypt" yaml:"bcrypt" ini:"bcrypt"`                                     // bcrypt 配置
	Argon2            Argon2Config  `json:"argon2" yaml:"argon2" ini:"argon2"`                                     // argon2 配置
	Argon2id          Argon2Config  `json:"argon2id" yaml:"argon2id" ini:"argon2id"`                               // argon2id 配置
	Scrypt            ScryptConfig  `json:"scrypt" yaml:"scrypt" ini:"scrypt"`                                     // scrypt 配置
	PBKDF2            PBKDF2Config  `json:"pbkdf2" yaml:"pbkdf2" ini:"pbkdf2"`                                     // pbkdf2 配置
}

// ensureDefaults 确保配置有合理的默认值
func (s *Security) ensureDefaults() {
	defaultCfg := defaultSecurityConfig()
	if s.PasswordAlgorithm == "" {
		s.PasswordAlgorithm = defaultCfg.PasswordAlgorithm
	}

	// 为每个算法配置设置默认值
	s.Bcrypt.ensureDefaults(defaultCfg.Bcrypt)
	s.Argon2.ensureDefaults(defaultCfg.Argon2)
	s.Argon2id.ensureDefaults(defaultCfg.Argon2id)
	s.Scrypt.ensureDefaults(defaultCfg.Scrypt)
	s.PBKDF2.ensureDefaults(defaultCfg.PBKDF2)
}

// BcryptConfig bcrypt 算法配置
type BcryptConfig struct {
	Cost int `json:"cost" yaml:"cost" ini:"cost"` // 计算成本，范围 4-31，默认 10
}

func (b *BcryptConfig) ensureDefaults(defaultCfg BcryptConfig) {
	if b.Cost == 0 {
		b.Cost = defaultCfg.Cost
	}
}

// Argon2Config argon2 算法配置
type Argon2Config struct {
	Memory      uint32 `json:"memory" yaml:"memory" ini:"memory"`                // 内存大小 (KB)
	Iterations  uint32 `json:"iterations" yaml:"iterations" ini:"iterations"`    // 迭代次数
	Parallelism uint8  `json:"parallelism" yaml:"parallelism" ini:"parallelism"` // 并行度
	SaltLength  uint32 `json:"salt_length" yaml:"salt_length" ini:"salt_length"` // 盐值长度
	KeyLength   uint32 `json:"key_length" yaml:"key_length" ini:"key_length"`    // 密钥长度
}

func (a *Argon2Config) ensureDefaults(defaultCfg Argon2Config) {
	if a.Memory == 0 {
		a.Memory = defaultCfg.Memory
	}
	if a.Iterations == 0 {
		a.Iterations = defaultCfg.Iterations
	}
	if a.Parallelism == 0 {
		a.Parallelism = defaultCfg.Parallelism
	}
	if a.SaltLength == 0 {
		a.SaltLength = defaultCfg.SaltLength
	}
	if a.KeyLength == 0 {
		a.KeyLength = defaultCfg.KeyLength
	}
}

// ScryptConfig scrypt 算法配置
type ScryptConfig struct {
	N       int `json:"n" yaml:"n" ini:"n"`                      // CPU/内存成本参数
	R       int `json:"r" yaml:"r" ini:"r"`                      // 块大小参数
	P       int `json:"p" yaml:"p" ini:"p"`                      // 并行度参数
	KeyLen  int `json:"key_len" yaml:"key_len" ini:"key_len"`    // 密钥长度
	SaltLen int `json:"salt_len" yaml:"salt_len" ini:"salt_len"` // 盐值长度
}

func (s *ScryptConfig) ensureDefaults(defaultCfg ScryptConfig) {
	if s.N == 0 {
		s.N = defaultCfg.N
	}
	if s.R == 0 {
		s.R = defaultCfg.R
	}
	if s.P == 0 {
		s.P = defaultCfg.P
	}
	if s.KeyLen == 0 {
		s.KeyLen = defaultCfg.KeyLen
	}
	if s.SaltLen == 0 {
		s.SaltLen = defaultCfg.SaltLen
	}
}

// PBKDF2Config PBKDF2 算法配置
type PBKDF2Config struct {
	Iterations int    `json:"iterations" yaml:"iterations" ini:"iterations"`    // 迭代次数
	KeyLength  int    `json:"key_length" yaml:"key_length" ini:"key_length"`    // 密钥长度
	HashFunc   string `json:"hash_func" yaml:"hash_func" ini:"hash_func"`       // 哈希函数 (sha1, sha256, sha512)
	SaltLength int    `json:"salt_length" yaml:"salt_length" ini:"salt_length"` // 盐值长度
}

func (p *PBKDF2Config) ensureDefaults(defaultCfg PBKDF2Config) {
	if p.Iterations == 0 {
		p.Iterations = defaultCfg.Iterations
	}
	if p.KeyLength == 0 {
		p.KeyLength = defaultCfg.KeyLength
	}
	if p.HashFunc == "" {
		p.HashFunc = defaultCfg.HashFunc
	}
	if p.SaltLength == 0 {
		p.SaltLength = defaultCfg.SaltLength
	}
}

var (
	ErrMismatchedHashAndPassword = errors.New("密码校验失败")     // ErrMismatchedHashAndPassword 密码不匹配错误
	ErrUnsupportedAlgorithm      = errors.New("不支持的hash算法") // ErrUnsupportedAlgorithm 不支持的算法错误
)

// Hasher 密码哈希器接口
type Hasher interface {
	Hash(password string) (string, error)          // Hash 对密码进行哈希
	Compare(password, hashedPassword string) error // Compare 比较密码和哈希值
}
