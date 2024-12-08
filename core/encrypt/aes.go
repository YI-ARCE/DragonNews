package encrypt

import (
	"bytes"
	enc "crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"strings"
	"unsafe"
)

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
func newECBEncrypter(b cipher.Block) cipher.BlockMode {
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
func newECBDecrypter(b cipher.Block) cipher.BlockMode {
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

type aes struct {
	cipher.Block
	key    string
	result []byte
}

// 填充数据
func padding(src string, blockSize int) []byte {
	padNum := blockSize - len(src)%blockSize
	pad := bytes.Repeat([]byte{byte(padNum)}, padNum)
	return append([]byte(src), pad...)
}

// 去掉填充数据
func unpadding(src string) []byte {
	n := len(src)
	if n == 0 {
		return []byte{}
	}
	unPadNum := int(src[n-1])
	return []byte(src)[:n-unPadNum]
}

func Aes(key string) *aes {
	lenght := len(key)
	if lenght < 16 {
		key += string(make([]byte, 16-lenght))
	}
	block, err := enc.NewCipher([]byte(key))
	if err != nil {
		panic(err.Error())
	}
	return &aes{block, key, []byte{}}
}

func (a *aes) Encrypt(cip string) *aes {
	// 加密
	dCip := padding(cip, a.Block.BlockSize())
	blockMode := newECBEncrypter(a.Block)
	blockMode.CryptBlocks(dCip, dCip)
	a.result = dCip
	return a
}

func (a *aes) Decrypt(cip string) *aes {
	var cipb []byte
	blockMode := newECBDecrypter(a.Block)
	cipb, err := hex.DecodeString(cip)
	if err != nil {
		cipb, err = base64.StdEncoding.DecodeString(cip)
		if err != nil {
			a.result = []byte{}
			return a
		}
	}
	blockMode.CryptBlocks(cipb, cipb)
	dCip := unpadding(string(cipb))
	a.result = dCip
	return a
}

func (a *aes) ToHex() string {
	return strings.ToUpper(hex.EncodeToString(a.result))
}

func (a *aes) ToBase64() string {
	return base64.StdEncoding.EncodeToString(a.result)
}

func (a *aes) ToString() string {
	return *(*string)(unsafe.Pointer(&a.result))
}
