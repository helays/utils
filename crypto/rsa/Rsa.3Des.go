package rsa

import (
	"crypto/cipher"
	"crypto/des"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"github.com/helays/utils/v2/crypto/sm4"
	"github.com/helays/utils/v2/logger/ulogs"
	"io/ioutil"
)

// //////////////////////////////Rsa加解密算法/////////////////////
func RsaEncrypt(origData []byte, filePth string) ([]byte, error) {
	PublicKey, err := ioutil.ReadFile(filePth)
	ulogs.Checkerr(err)
	block, _ := pem.Decode(PublicKey)
	if block == nil {
		return []byte{}, errors.New("public key empty")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	ulogs.Checkerr(err)
	pub := pubInterface.(*rsa.PublicKey)
	//print(pub.Size())
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)

}

func RsaDecrypt(ciphertext []byte, filePth string) ([]byte, error) {
	privateKey, err := ioutil.ReadFile(filePth)
	block, _ := pem.Decode(privateKey)

	//解析PKCS1格式的私钥
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 解密
	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}

// ////////////////////////提供3Des加解密的方法/////////////////////////////////
func ThriDESDeCrypt(crypted, key []byte) []byte {

	block, _ := des.NewTripleDESCipher(key[:24])
	//创建切片
	context := make([]byte, len(crypted))
	//设置解密方式
	blockMode := cipher.NewCBCDecrypter(block, key[24:])
	//解密密文到数组
	blockMode.CryptBlocks(context, crypted)
	//去补码

	context = sm4.Pkcs5UnPadding(context)
	return context
}

// ThriDESEnCrypt 加密
func ThriDESEnCrypt(origData, key []byte) []byte {
	//获取block块
	block, err := des.NewTripleDESCipher(key[:24])
	//补码
	ulogs.Checkerr(err)
	origData = sm4.Pkcs5Padding(origData, block.BlockSize())
	//设置加密方式为 3DES  使用3条56位的密钥对数据进行三次加密

	blockMode := cipher.NewCBCEncrypter(block, key[24:])

	//创建明文长度的数组
	crypted := make([]byte, len(origData))

	//加密明文
	blockMode.CryptBlocks(crypted, origData)
	return crypted

}
