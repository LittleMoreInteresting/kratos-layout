package aes_cbc

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

//加密过程：
//  1、处理数据，对数据进行填充，采用PKCS7（当密钥长度不够时，缺几位补几个几）的方式。
//  2、对数据进行加密，采用AES加密方法中CBC加密模式
//  3、对得到的加密数据，进行base64加密，得到字符串
// 解密过程相反

//16,24,32位字符串的话，分别对应AES-128，AES-192，AES-256 加密方法
//key
type AesCBC struct {
	key []byte
}
type Option func(aes *AesCBC)

func WithKey(key string) Option {
	return func(aes *AesCBC) {
		aes.key = []byte(key)
	}
}
func NewAesCBC(opt ...Option) *AesCBC {
	aesObj := &AesCBC{}
	for _, o := range opt {
		o(aesObj)
	}
	return aesObj
}

//pkcs7Padding
func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

//pkcs7UnPadding
func pkcs7UnPadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("加密字符串错误！")
	}
	unPadding := int(data[length-1])
	return data[:(length - unPadding)], nil
}

//AesEncrypt 加密
func AesEncrypt(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	encryptBytes := pkcs7Padding(data, blockSize)
	crypted := make([]byte, len(encryptBytes))
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	blockMode.CryptBlocks(crypted, encryptBytes)
	return crypted, nil
}

//AesDecrypt 解密
func AesDecrypt(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	crypted := make([]byte, len(data))
	blockMode.CryptBlocks(crypted, data)
	crypted, err = pkcs7UnPadding(crypted)
	if err != nil {
		return nil, err
	}
	return crypted, nil
}

//EncryptByAes Aes加密 后 base64
func (ac *AesCBC) EncryptByAes(data []byte) (string, error) {
	res, err := AesEncrypt(data, ac.key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(res), nil
}

//DecryptByAes base64解码后 Aes 解密
func (ac *AesCBC) DecryptByAes(data string) ([]byte, error) {
	dataByte, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	return AesDecrypt(dataByte, ac.key)
}

// EncryptyPhone
func (ac *AesCBC) EncryptPhone(phone string) (res string) {
	if len(phone) <= 11 {
		res, _ = ac.EncryptByAes([]byte(phone))
		return
	}
	return phone
}

// DecryptPhone
func (ac *AesCBC) DecryptPhone(phone string) (res string) {
	if len(phone) <= 11 {
		return phone
	}
	p, _ := ac.DecryptByAes(phone)
	return string(p)
}
