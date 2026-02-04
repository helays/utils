package jwtkit

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/helays/utils/v2/close/vclose"
	"github.com/helays/utils/v2/tools"
)

var (
	ErrUnsupportedAlgorithm = errors.New("不支持的签名算法")
	ErrKeyFileNotFound      = errors.New("密钥文件不存在")
	ErrInvalidKeyFormat     = errors.New("无效的密钥格式")
)

type keyLoader struct {
	config     *Options
	signingKey []byte // 签名密钥
	verifyKey  []byte // 验证密钥
}

func newKeyLoader(cfg *Options) *keyLoader {
	return &keyLoader{
		config: cfg,
	}
}

func (l *keyLoader) setSigningKey() error {
	switch l.config.Algorithm {
	case HS256, HS384, HS512:
		return l.loadHMACKey()
	case RS256, RS384, RS512:
		return l.loadRSAKey()
	case PS256, PS384, PS512:
		return l.loadPSSKey()
	case ES256, ES384, ES512:
		return l.loadECDSAKey()
	case EdDSA:
		return l.loadEdDSAKey()
	case SigningMethodNone:
		return nil
	default:
		return ErrUnsupportedAlgorithm
	}
}

func (l *keyLoader) mkdir() error {
	return tools.Mkdir(l.config.KeyStoreDir)
}

func (l *keyLoader) loadHMACKey() error {
	if l.config.HMAC.Secret != "" {
		l.signingKey = []byte(l.config.HMAC.Secret)
		l.verifyKey = l.signingKey
		return nil
	}
	if l.config.HMAC.SecretPath != "" {
		p := tools.Fileabs(l.config.HMAC.SecretPath)
		b, err := tools.FileGetContents(p)
		if err != nil {
			return err
		}
		l.signingKey = b
		l.verifyKey = l.signingKey
		return nil
	}
	if !l.config.AutoGenerate {
		return errors.New("未配置HMAC密钥")
	}
	if err := l.mkdir(); err != nil {
		return err
	}
	p := filepath.Join(l.config.KeyStoreDir, "hmac.key")
	f, err := os.OpenFile(p, os.O_CREATE|os.O_RDWR, 0600)
	defer vclose.Close(f)
	if err != nil {
		return err
	}
	b, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	if len(b) > 0 {
		l.signingKey = b
		l.verifyKey = l.signingKey
		return nil
	}
	key := make([]byte, l.getHMACKeySize())
	if _, err = rand.Read(key); err != nil {
		return fmt.Errorf("无法生成HMAC密钥: %v", err)
	}
	if _, err = f.Write(key); err != nil {
		return fmt.Errorf("无法保存HMAC密钥: %v", err)
	}
	l.signingKey = key
	l.verifyKey = l.signingKey
	return nil
}

// getHMACKeySize 根据算法确定密钥长度
func (l *keyLoader) getHMACKeySize() int {
	switch l.config.Algorithm {
	case HS256:
		return 32 // 256位
	case HS384:
		return 48 // 384位
	case HS512:
		return 64 // 512位
	default:
		return 32 // 默认256位
	}
}

func (l *keyLoader) loadRSAKey() error {
	if l.config.RSA.Private != "" {
		l.signingKey = []byte(l.config.RSA.Private)
		l.verifyKey = []byte(l.config.RSA.Public)
		return nil
	}
	if l.config.RSA.PrivatePath != "" {
		privatePath := tools.Fileabs(l.config.RSA.PrivatePath)
		publicPath := tools.Fileabs(l.config.RSA.PublicPath)
		private, err := tools.FileGetContents(privatePath)
		if err != nil {
			return err
		}
		public, err := tools.FileGetContents(publicPath)
		if err != nil {
			return err
		}
		l.signingKey = private
		l.verifyKey = public
		return nil
	}
	if !l.config.AutoGenerate {
		return errors.New("未配置RSA密钥")
	}

	privatePath := filepath.Join(l.config.KeyStoreDir, "rsa.key")
	publicPath := filepath.Join(l.config.KeyStoreDir, "rsa.pub")
	// 尝试加载现有密钥
	if l.tryLoadExistingKeyPair(privatePath, publicPath) {
		return nil
	}
	// 生成新的RSA密钥对
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("生成RSA密钥失败: %v", err)
	}
	// 保存私钥
	privatePEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	if err = l.saveKeyToFile(privatePath, privatePEM); err != nil {
		return fmt.Errorf("保存RSA私钥失败: %v", err)
	}
	// 保存公钥
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return fmt.Errorf("序列化RSA公钥失败: %v", err)
	}
	publicPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})
	if err = l.saveKeyToFile(publicPath, publicPEM); err != nil {
		return fmt.Errorf("保存RSA公钥失败: %v", err)
	}
	l.signingKey = privatePEM
	l.verifyKey = publicPEM
	return nil
}

func (l *keyLoader) loadPSSKey() error {
	// PSS使用RSA密钥，复用RSA的加载逻辑
	return l.loadRSAKey()
}

func (l *keyLoader) loadECDSAKey() error {
	if l.config.ECDSA.Private != "" {
		l.signingKey = []byte(l.config.ECDSA.Private)
		l.verifyKey = []byte(l.config.ECDSA.Public)
		return nil
	}
	if l.config.ECDSA.PrivatePath != "" {
		privatePath := tools.Fileabs(l.config.ECDSA.PrivatePath)
		publicPath := tools.Fileabs(l.config.ECDSA.PublicPath)
		private, err := tools.FileGetContents(privatePath)
		if err != nil {
			return err
		}
		public, err := tools.FileGetContents(publicPath)
		if err != nil {
			return err
		}
		l.signingKey = private
		l.verifyKey = public
		return nil
	}
	if !l.config.AutoGenerate {
		return errors.New("未配置ECDSA密钥")
	}
	privatePath := filepath.Join(l.config.KeyStoreDir, "ecdsa.key")
	publicPath := filepath.Join(l.config.KeyStoreDir, "ecdsa.pub")
	// 尝试加载现有密钥
	if l.tryLoadExistingKeyPair(privatePath, publicPath) {
		return nil
	}
	// 根据算法选择椭圆曲线
	var curve elliptic.Curve
	switch l.config.Algorithm {
	case ES256:
		curve = elliptic.P256()
	case ES384:
		curve = elliptic.P384()
	case ES512:
		curve = elliptic.P521()
	default:
		curve = elliptic.P256() // 默认使用P256
	}

	// 生成新的ECDSA密钥对
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return fmt.Errorf("生成ECDSA密钥失败: %v", err)
	}

	// 保存私钥
	privateBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return fmt.Errorf("序列化ECDSA私钥失败: %v", err)
	}
	privatePEM := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privateBytes,
	})
	if err = l.saveKeyToFile(privatePath, privatePEM); err != nil {
		return fmt.Errorf("保存ECDSA私钥失败: %v", err)
	}

	// 保存公钥
	publicBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return fmt.Errorf("序列化ECDSA公钥失败: %v", err)
	}
	publicPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PUBLIC KEY",
		Bytes: publicBytes,
	})
	if err = l.saveKeyToFile(publicPath, publicPEM); err != nil {
		return fmt.Errorf("保存ECDSA公钥失败: %v", err)
	}

	l.signingKey = privatePEM
	l.verifyKey = publicPEM
	return nil
}

func (l *keyLoader) loadEdDSAKey() error {
	if l.config.EdDSA.Private != "" {
		l.signingKey = []byte(l.config.EdDSA.Private)
		l.verifyKey = []byte(l.config.EdDSA.Public)
		return nil
	}
	if l.config.EdDSA.PrivatePath != "" {
		privatePath := tools.Fileabs(l.config.EdDSA.PrivatePath)
		publicPath := tools.Fileabs(l.config.EdDSA.PublicPath)
		private, err := tools.FileGetContents(privatePath)
		if err != nil {
			return err
		}
		public, err := tools.FileGetContents(publicPath)
		if err != nil {
			return err
		}
		l.signingKey = private
		l.verifyKey = public
		return nil
	}
	if !l.config.AutoGenerate {
		return errors.New("未配置EdDSA密钥")
	}
	privatePath := filepath.Join(l.config.KeyStoreDir, "ed25519.key")
	publicPath := filepath.Join(l.config.KeyStoreDir, "ed25519.pub")

	// 尝试加载现有密钥
	if l.tryLoadExistingKeyPair(privatePath, publicPath) {
		return nil
	}

	// 生成新的Ed25519密钥对
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return fmt.Errorf("生成EdDSA密钥失败: %v", err)
	}

	// 保存私钥
	privateBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return fmt.Errorf("序列化EdDSA私钥失败: %v", err)
	}
	privatePEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateBytes,
	})
	if err = l.saveKeyToFile(privatePath, privatePEM); err != nil {
		return fmt.Errorf("保存EdDSA私钥失败: %v", err)
	}

	// 保存公钥
	publicBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return fmt.Errorf("序列化EdDSA公钥失败: %v", err)
	}
	publicPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicBytes,
	})
	if err = l.saveKeyToFile(publicPath, publicPEM); err != nil {
		return fmt.Errorf("保存EdDSA公钥失败: %v", err)
	}

	l.signingKey = privatePEM
	l.verifyKey = publicPEM
	return nil
}

// tryLoadExistingKeyPair 尝试加载现有的密钥对
func (l *keyLoader) tryLoadExistingKeyPair(privatePath, publicPath string) bool {
	private, err := tools.FileGetContents(privatePath)
	if err != nil {
		return false
	}
	public, err := tools.FileGetContents(publicPath)
	if err != nil {
		return false
	}

	if len(private) > 0 && len(public) > 0 {
		l.signingKey = private
		l.verifyKey = public
		return true
	}
	return false
}

// saveKeyToFile 保存密钥到文件
func (l *keyLoader) saveKeyToFile(path string, data []byte) error {
	// 确保目录存在
	dir := filepath.Dir(path)
	if err := tools.Mkdir(dir); err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer vclose.Close(file)

	_, err = file.Write(data)
	return err
}
