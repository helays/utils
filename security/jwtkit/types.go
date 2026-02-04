// Package jwtkit
// 支持设置多种签名算法
// 默认使用HMAC
// 支持RSA、ECDSA、HMAC
// 支持手动设置密钥内容，密钥路径
// 支持自动生成密钥
package jwtkit

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/helays/utils/v2/dataType"
)

type Options struct {
	Algorithm    Algorithm      `json:"algorithm" yaml:"algorithm" ini:"algorithm"`             // 签名算法
	AutoGenerate bool           `json:"auto_generate" yaml:"auto_generate" ini:"auto_generate"` // 是否自动生成密钥
	KeyStoreDir  string         `json:"key_store_dir" yaml:"key_store_dir" ini:"key_store_dir"` // 密钥存储目录
	HMAC         HMAC           `json:"hmac" yaml:"hmac" ini:"hmac"`                            // HMAC
	RSA          AsymmetricKeys `json:"rsa" yaml:"rsa" ini:"rsa"`                               // RSA
	ECDSA        AsymmetricKeys `json:"ecdsa" yaml:"ecdsa" ini:"ecdsa"`                         // ECDSA
	EdDSA        AsymmetricKeys `json:"eddsa" yaml:"eddsa" ini:"eddsa"`                         // EdDSA
}

type HMAC struct {
	Secret     string `json:"secret" yaml:"secret" ini:"secret"`                // HMAC密钥
	SecretPath string `json:"secret_path" yaml:"secret_path" ini:"secret_path"` // HMAC密钥路径
}

// AsymmetricKeys 非对称密钥
type AsymmetricKeys struct {
	Private     string `json:"private" yaml:"private" ini:"private"`
	Public      string `json:"public" yaml:"public" ini:"public"`
	PrivatePath string `json:"private_path" yaml:"private_path" ini:"private_path"`
	PublicPath  string `json:"public_path" yaml:"public_path" ini:"public_path"`
}

// StandardClaims 标准Claims
// exp 过期时间
// nbf 生效时间
// iat 签发时间
// iss 签发人 签发者URL或者应用标识
// aud 接收人 给谁使用的 也可以是URL或者应用标识
// sub 主题
// jti 唯一标识
type StandardClaims struct {
	UserId    dataType.IntString[int64] `json:"user_id"`
	TenantID  dataType.IntString[int64] `json:"tenant_id"`  // 租户ID
	Timestamp time.Time                 `json:"timestamp"`  // 登录时间
	LoginIp   string                    `json:"login_ip"`   // 登录IP
	ExpiresAt *time.Time                `json:"expires_at"` // 过期时间
	Extras    map[string]any            `json:"extras"`     // 扩展字段
	jwt.RegisteredClaims
}
