package encrypt

import (
	"bytes"
	enc "crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
)

// aes AES加密结构
// 包含加密块、密钥和结果

type aes struct {
	cipher.Block
	key    string
	result []byte
}

// padding 填充数据到指定块大小
//
// 参数：
//   - src: 原始字符串
//   - blockSize: 块大小
//
// 返回值：
//   - []byte: 填充后的数据
func padding(src string, blockSize int) []byte {
	padNum := blockSize - len(src)%blockSize
	pad := bytes.Repeat([]byte{byte(padNum)}, padNum)
	return append([]byte(src), pad...)
}

// unpadding 去掉填充数据
//
// 参数：
//   - src: 填充后的数据
//
// 返回值：
//   - []byte: 原始数据
func unpadding(src []byte) []byte {
	n := len(src)
	if n == 0 {
		return []byte{}
	}
	unPadNum := int(src[n-1])
	return src[:n-unPadNum]
}

// Aes 创建AES加密对象
//
// 参数：
//   - key: 加密密钥
//
// 返回值：
//   - *aes: AES加密对象
//   - error: 错误信息
//
// 说明：
//   - 密钥长度不足16字节时，会自动填充到16字节
//   - 密钥长度超过32字节时，会截断到32字节
func Aes(key string) (*aes, error) {
	length := len(key)
	if length < 16 {
		key += string(make([]byte, 16-length))
	} else if length > 32 {
		key = key[:32]
	}
	block, err := enc.NewCipher([]byte(key))
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}
	return &aes{block, key, []byte{}}, nil
}

// Encrypt 加密字符串
//
// 参数：
//   - cip: 待加密的字符串
//
// 返回值：
//   - *aes: AES加密对象
//   - error: 错误信息
//
// 说明：
//   - 使用CBC模式加密
//   - 自动生成随机IV
//   - 将IV添加到加密结果前面
func (a *aes) Encrypt(cip string) (*aes, error) {
	// 加密
	dCip := padding(cip, a.Block.BlockSize())

	// 生成随机IV
	iv := make([]byte, a.Block.BlockSize())
	_, err := rand.Read(iv)
	if err != nil {
		return a, fmt.Errorf("failed to generate IV: %w", err)
	}

	// 使用CBC模式
	blockMode := cipher.NewCBCEncrypter(a.Block, iv)
	blockMode.CryptBlocks(dCip, dCip)

	// 将IV添加到加密结果前面
	a.result = append(iv, dCip...)
	return a, nil
}

// Decrypt 解密字符串
//
// 参数：
//   - cip: 待解密的字符串（支持hex或base64格式）
//
// 返回值：
//   - *aes: AES加密对象
//   - error: 错误信息
//
// 说明：
//   - 自动检测并处理hex或base64格式的密文
//   - 从密文中提取IV
//   - 使用CBC模式解密
//   - 自动去除填充
func (a *aes) Decrypt(cip string) (*aes, error) {
	var cipb []byte
	var err error

	cipb, err = hex.DecodeString(cip)
	if err != nil {
		cipb, err = base64.StdEncoding.DecodeString(cip)
		if err != nil {
			a.result = []byte{}
			return a, fmt.Errorf("failed to decode ciphertext: %w", err)
		}
	}

	// 检查数据长度
	blockSize := a.Block.BlockSize()
	if len(cipb) < blockSize {
		a.result = []byte{}
		return a, fmt.Errorf("ciphertext too short")
	}

	// 提取IV
	iv := cipb[:blockSize]
	cipb = cipb[blockSize:]

	// 使用CBC模式
	blockMode := cipher.NewCBCDecrypter(a.Block, iv)
	blockMode.CryptBlocks(cipb, cipb)

	// 去填充
	dCip := unpadding(cipb)
	a.result = dCip
	return a, nil
}

// ToHex 将结果转换为十六进制字符串
//
// 返回值：
//   - string: 十六进制字符串（大写）
func (a *aes) ToHex() string {
	return strings.ToUpper(hex.EncodeToString(a.result))
}

// ToBase64 将结果转换为Base64字符串
//
// 返回值：
//   - string: Base64字符串
func (a *aes) ToBase64() string {
	return base64.StdEncoding.EncodeToString(a.result)
}

// ToString 将结果转换为字符串
//
// 返回值：
//   - string: 字符串
func (a *aes) ToString() string {
	return string(a.result)
}
