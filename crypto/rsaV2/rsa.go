package rsaV2

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
)

const (
	privateBlockType = "RSA PRIVATE KEY"
	publicBlockType  = "RSA PUBLIC KEY"
)

func decodePublicKey(publicKey []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("failed to decode PEM block containing private key")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pub.(*rsa.PublicKey), err
}

func decodePrivateKey(privateKey []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil || block.Type != privateBlockType {
		return nil, errors.New("failed to decode PEM block containing private key")
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func split(buf []byte, lim int) [][]byte {
	var chunk []byte
	chunks := make([][]byte, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf[:len(buf)])
	}
	return chunks
}

// GenRsaPrivateKey 生成私钥
func GenRsaPrivateKey(bits int) (*rsa.PrivateKey, []byte, error) {
	key, err := rsa.GenerateKey(rand.Reader, bits) //生成一对具有指定字位数的RSA密钥
	if err != nil {
		return nil, nil, err
	}
	if err = key.Validate(); err != nil {
		return nil, nil, err
	}
	return key, pem.EncodeToMemory(&pem.Block{Type: privateBlockType, Bytes: x509.MarshalPKCS1PrivateKey(key)}), nil
}

// GenRsaPublicKey 获取公钥
func GenRsaPublicKey(privateKey []byte) ([]byte, error) {
	key, err := decodePrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	_byte, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(&pem.Block{
		Type:  publicBlockType,
		Bytes: _byte,
	}), nil
}

// GenRsaPriPubKey 生成公钥私钥
func GenRsaPriPubKey(bits int) ([]byte, []byte, error) {
	key, prikey, err := GenRsaPrivateKey(bits)
	if err != nil {
		return nil, nil, err
	}
	_byte, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		return nil, nil, err
	}
	pubKey := pem.EncodeToMemory(&pem.Block{
		Type:  publicBlockType,
		Bytes: _byte,
	})
	return prikey, pubKey, nil
}

// EncryptWithPublicKey
// 公钥加密
func EncryptWithPublicKey(publicKey []byte, bs string) (string, error) {
	enc, err := EncryptBytWithPublicKey(publicKey, []byte(bs))
	if err != nil {
		return "", err
	}
	return string(enc), nil
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

// EncryptBytWithPublicKey
// 公钥加密，传入byte
// 可加密任意长度数据
func EncryptBytWithPublicKey(publicKey []byte, bs []byte) ([]byte, error) {
	key, err := decodePublicKey(publicKey)
	if err != nil {
		return nil, err
	}

	partLen := key.N.BitLen()/8 - 42
	// 根据partLen 分割
	var buffer bytes.Buffer
	for _, chunk := range split(bs, partLen) {
		encryptBytes, err := rsa.EncryptPKCS1v15(rand.Reader, key, chunk)
		if err != nil {
			return nil, err
		}
		buffer.Write(encryptBytes)
	}
	buf := make([]byte, base64.StdEncoding.EncodedLen(buffer.Len()))
	base64.StdEncoding.Encode(buf, buffer.Bytes())
	return buf, nil
}

func DecryptBytWithPrivateKey(privateKey []byte, encryptData []byte) ([]byte, error) {
	key, err := decodePrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	dbuf := make([]byte, base64.StdEncoding.DecodedLen(len(encryptData)))
	n, err := base64.StdEncoding.Decode(dbuf, encryptData)
	if err != nil {
		return nil, err
	}
	encryptBytes := dbuf[:n]
	partLen := key.N.BitLen() / 8
	var buffer bytes.Buffer
	for _, chunk := range split(encryptBytes, partLen) {
		byt, err := rsa.DecryptPKCS1v15(rand.Reader, key, chunk)
		if err != nil {
			return nil, err
		}
		buffer.Write(byt)
	}
	return buffer.Bytes(), nil
}
