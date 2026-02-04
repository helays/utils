package jwtkit

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"helay.net/go/utils/v3/tools"
	"helay.net/go/utils/v3/tools/sonyflakekit"
)

type JWTManager struct {
	config    *Options
	method    jwt.SigningMethod
	keyLoader *keyLoader
	snowflake *sonyflakekit.IDGenerator
}

func NewJWTManager(config *Options, snowflake *sonyflakekit.IDGenerator) (*JWTManager, error) {
	config.KeyStoreDir = tools.Fileabs(tools.Ternary(config.KeyStoreDir == "", "keys", config.KeyStoreDir))
	kl := newKeyLoader(config)
	if err := kl.setSigningKey(); err != nil {
		return nil, err
	}
	m := &JWTManager{config: config, keyLoader: kl, snowflake: snowflake}
	method, err := m.getSigningMethod(config.Algorithm)
	if err != nil {
		return nil, err
	}
	m.method = method
	return m, nil

}

func (m *JWTManager) getSigningMethod(alg Algorithm) (jwt.SigningMethod, error) {
	switch alg {
	case HS256:
		return jwt.SigningMethodHS256, nil
	case HS384:
		return jwt.SigningMethodHS384, nil
	case HS512:
		return jwt.SigningMethodHS512, nil
	case RS256:
		return jwt.SigningMethodRS256, nil
	case RS384:
		return jwt.SigningMethodRS384, nil
	case RS512:
		return jwt.SigningMethodRS512, nil
	case PS256:
		return jwt.SigningMethodPS256, nil
	case PS384:
		return jwt.SigningMethodPS384, nil
	case PS512:
		return jwt.SigningMethodPS512, nil
	case ES256:
		return jwt.SigningMethodES256, nil
	case ES384:
		return jwt.SigningMethodES384, nil
	case ES512:
		return jwt.SigningMethodES512, nil
	case EdDSA:
		return jwt.SigningMethodEdDSA, nil
	case SigningMethodNone:
		return jwt.SigningMethodNone, nil
	}
	return nil, ErrUnsupportedAlgorithm
}

func (m *JWTManager) GenerateToken(claims *StandardClaims) (string, error) {
	jit, err := m.snowflake.GenerateID()
	if err != nil {
		return "", err
	}
	claims.ID = tools.Any2string(jit)
	claims.Timestamp = time.Now()
	token := jwt.NewWithClaims(m.method, claims)
	return token.SignedString(m.keyLoader.signingKey)
}

func (m *JWTManager) ValidateToken(tokenString string) (*StandardClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法是否匹配
		if token.Method != m.method {
			return nil, jwt.ErrSignatureInvalid
		}
		return m.keyLoader.verifyKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*StandardClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrInvalidKey
}
