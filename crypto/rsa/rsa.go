package rsa

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"

	"github.com/helays/utils/v2/config"
	"github.com/helays/utils/v2/crypto/sha256"
)

// RsaVerify 签名验证，用公钥进行验证
// msg 源实内容
// 签名
func RsaVerify(msg string, _sign string) error {
	// 数据hash
	msgHashbyt, err := sha256.Sha256(msg)
	if err != nil {
		return err
	}
	// key
	block, _ := pem.Decode(config.PublicKeyByt)
	if block == nil {
		return errors.New("public key error")
	}
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}
	sign, err := base64.StdEncoding.DecodeString(_sign)
	if err != nil {
		return nil
	}
	return rsa.VerifyPKCS1v15(pubKey.(*rsa.PublicKey), crypto.SHA256, msgHashbyt, sign)
}

// RsaSign 内容签名，用私钥签名
func RsaSign(msg string) (string, error) {
	msgHashByt, err := sha256.Sha256(msg)
	if err != nil {
		return "", errors.New("Sha256(msg)---->" + err.Error())
	}
	block, _ := pem.Decode(config.PrivateKeyByt)
	if block == nil {
		return "", errors.New("Private key error")
	}
	var priKey interface{}
	switch block.Type {
	case "RSA PRIVATE KEY":
		priKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	case "PRIVATE KEY":
		priKey, err = x509.ParsePKCS8PrivateKey(block.Bytes)
	}
	//priKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", errors.New("block.Type：" + block.Type + "--->" + err.Error())
	}
	_sig, err := rsa.SignPKCS1v15(rand.Reader, priKey.(*rsa.PrivateKey), crypto.SHA256, msgHashByt)
	if err != nil {
		return "", errors.New("rsa.SignPKCS1v15(rand.Reader, priKey, crypto.SHA256, msgHashByt)--->" + err.Error())
	}
	return base64.StdEncoding.EncodeToString(_sig), nil
}

const (
	privateBlockType = "RSA PRIVATE KEY"
	publicBlockType  = "RSA PUBLIC KEY"
)

func decodePublicKey(publicKey []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("failed to decode PEM block containing private key")
	}
	return x509.ParsePKCS1PublicKey(block.Bytes)
}

func decodePrivateKey(privateKey []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil || block.Type != privateBlockType {
		return nil, errors.New("failed to decode PEM block containing private key")
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

// EncryptWithPublicKey 公钥加密
func EncryptWithPublicKey(publicKey []byte, bs string) (string, error) {
	key, err := decodePublicKey(publicKey)
	if err != nil {
		return "", err
	}
	encryptBytes, err := rsa.EncryptPKCS1v15(rand.Reader, key, []byte(bs))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encryptBytes), nil
}

// DecryptWithPrivateKey 私钥解密
func DecryptWithPrivateKey(privateKey []byte, encryptData string) ([]byte, error) {
	key, err := decodePrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	encryptBytes, err := base64.StdEncoding.DecodeString(encryptData)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, key, encryptBytes)
}
