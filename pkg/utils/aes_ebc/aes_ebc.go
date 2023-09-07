package aes_ebc

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

type AesEBC struct {
	key []byte
}

type Option func(aes *AesEBC)

func WithKey(key string) Option {
	return func(aes *AesEBC) {
		aes.key = []byte(key)
	}
}

func NewAesEBC(opt ...Option) *AesEBC {
	aesObj := &AesEBC{}
	for _, o := range opt {
		o(aesObj)
	}
	return aesObj
}

func (ac *AesEBC) Encrypt(in string) string {
	block, err := aes.NewCipher(ac.key)
	if err != nil {
		return in
	}
	blockSize := block.BlockSize()
	encryptBytes := PKCS5Padding([]byte(in), blockSize)
	crypted := make([]byte, len(encryptBytes))
	blockMode := NewECBEncrypter(block)
	blockMode.CryptBlocks(crypted, encryptBytes)
	return base64.StdEncoding.EncodeToString(crypted)
}

func (ac *AesEBC) Decrypt(in string) string {
	dataByte, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		panic(err)
	}
	block, err := aes.NewCipher(ac.key)
	if err != nil {
		panic(err)
	}
	blockMode := NewECBDecrypter(block)
	orig := make([]byte, len(dataByte))
	blockMode.CryptBlocks(orig, dataByte)
	orig = PKCS5UnPadding(orig)
	return string(orig)
}

func (ac *AesEBC) EncryptPhone(phone string) (res string) {
	if len(phone) <= 11 {
		res = ac.Encrypt(phone)
		return
	}
	return phone
}

// DecryptPhone
func (ac *AesEBC) DecryptPhone(phone string) (res string) {
	if len(phone) <= 11 {
		return phone
	}
	p := ac.Decrypt(phone)
	return p
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

type ecb struct {
	b         cipher.Block
	blockSize int
}

func newECB(b cipher.Block) *ecb {
	return &ecb{
		b:         b,
		blockSize: b.BlockSize(),
	}
}

type ecbEncrypter ecb

// NewECBEncrypter returns a BlockMode which encrypts in electronic code book
// mode, using the given Block.
func NewECBEncrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbEncrypter)(newECB(b))
}
func (x *ecbEncrypter) BlockSize() int { return x.blockSize }
func (x *ecbEncrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Encrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

type ecbDecrypter ecb

// NewECBDecrypter returns a BlockMode which decrypts in electronic code book
// mode, using the given Block.
func NewECBDecrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbDecrypter)(newECB(b))
}
func (x *ecbDecrypter) BlockSize() int { return x.blockSize }
func (x *ecbDecrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Decrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}
