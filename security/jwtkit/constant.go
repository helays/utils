package jwtkit

type Algorithm string

func (a Algorithm) String() string {
	return string(a)
}

// HMAC 系列
const (
	HS256 Algorithm = "HS256"
	HS384 Algorithm = "HS384"
	HS512 Algorithm = "HS512"
)

// RSA 签名系列
const (
	RS256 Algorithm = "RS256"
	RS384 Algorithm = "RS384"
	RS512 Algorithm = "RS512"
)

// PSS 签名系列
const (
	PS256 Algorithm = "PS256"
	PS384 Algorithm = "PS384"
	PS512 Algorithm = "PS512"
)

// ECDSA 签名系列
const (
	ES256 Algorithm = "ES256"
	ES384 Algorithm = "ES384"
	ES512 Algorithm = "ES512"
)

const EdDSA Algorithm = "EdDSA" // EdDSA 签名系列

const SigningMethodNone = "none" // 无签名
