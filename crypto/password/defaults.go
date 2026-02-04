package password

// DefaultSecurityConfig 返回默认的安全配置
func defaultSecurityConfig() *Security {
	return &Security{
		PasswordAlgorithm: HashBcrypt,
		Bcrypt:            BcryptConfig{Cost: 10},
		Argon2: Argon2Config{
			Memory:      64 * 1024, // 64MB
			Iterations:  3,
			Parallelism: 4,
			SaltLength:  16,
			KeyLength:   32,
		},
		Argon2id: Argon2Config{
			Memory:      64 * 1024, // 64MB
			Iterations:  1,
			Parallelism: 4,
			SaltLength:  16,
			KeyLength:   32,
		},
		Scrypt: ScryptConfig{
			N:       32768,
			R:       8,
			P:       1,
			KeyLen:  32,
			SaltLen: 16,
		},
		PBKDF2: PBKDF2Config{
			Iterations: 100000,
			KeyLength:  32,
			HashFunc:   "sha256",
			SaltLength: 16,
		},
	}
}
