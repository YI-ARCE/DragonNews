package test

import (
	"testing"
	"yiarce/core/encrypt"
)

// TestAesCreation 测试AES加密对象创建功能
func TestAesCreation(t *testing.T) {
	// 测试创建AES加密对象
	key := "test_key_1234567" // 16字节
	aesObj, err := encrypt.Aes(key)
	if err != nil {
		t.Errorf("Aes creation error: %v", err)
	}
	if aesObj == nil {
		t.Error("Aes creation returned nil")
	}

	// 测试创建AES加密对象的边界条件
	// 测试key长度 < 16的情况
	shortKey := "test_key"
	aesObj, err = encrypt.Aes(shortKey)
	if err != nil {
		t.Errorf("Aes creation with short key error: %v", err)
	}
	if aesObj == nil {
		t.Error("Aes creation with short key returned nil")
	}

	// 测试key长度 > 32的情况
	longKey := "test_key_1234567890123456_test_key_1234567890"
	aesObj, err = encrypt.Aes(longKey)
	if err != nil {
		t.Errorf("Aes creation with long key error: %v", err)
	}
	if aesObj == nil {
		t.Error("Aes creation with long key returned nil")
	}
}

// TestAesEncryptDecrypt 测试AES加密和解密功能
func TestAesEncryptDecrypt(t *testing.T) {
	// 测试数据
	key := "test_key_1234567" // 16字节
	plaintext := "Hello, World!"

	// 创建AES加密对象
	aesObj, err := encrypt.Aes(key)
	if err != nil {
		t.Errorf("Aes creation error: %v", err)
		return
	}

	// 加密
	encrypted, err := aesObj.Encrypt(plaintext)
	if err != nil {
		t.Errorf("Aes encrypt error: %v", err)
		return
	}

	// 将加密结果转换为十六进制
	hexEncrypted := encrypted.ToHex()
	if hexEncrypted == "" {
		t.Error("Aes encrypt returned empty result")
		return
	}

	// 解密
	decrypted, err := aesObj.Decrypt(hexEncrypted)
	if err != nil {
		t.Errorf("Aes decrypt error: %v", err)
		return
	}

	// 将解密结果转换为字符串
	decryptedText := decrypted.ToString()
	if decryptedText != plaintext {
		t.Errorf("Aes decrypt returned wrong result, expected: %s, got: %s", plaintext, decryptedText)
	}
}

// TestAesBase64 测试AES加密和解密功能（Base64格式）
func TestAesBase64(t *testing.T) {
	// 测试数据
	key := "test_key_1234567" // 16字节
	plaintext := "Hello, Base64!"

	// 创建AES加密对象
	aesObj, err := encrypt.Aes(key)
	if err != nil {
		t.Errorf("Aes creation error: %v", err)
		return
	}

	// 加密
	encrypted, err := aesObj.Encrypt(plaintext)
	if err != nil {
		t.Errorf("Aes encrypt error: %v", err)
		return
	}

	// 将加密结果转换为Base64
	base64Encrypted := encrypted.ToBase64()
	if base64Encrypted == "" {
		t.Error("Aes encrypt returned empty result")
		return
	}

	// 解密
	decrypted, err := aesObj.Decrypt(base64Encrypted)
	if err != nil {
		t.Errorf("Aes decrypt error: %v", err)
		return
	}

	// 将解密结果转换为字符串
	decryptedText := decrypted.ToString()
	if decryptedText != plaintext {
		t.Errorf("Aes decrypt returned wrong result, expected: %s, got: %s", plaintext, decryptedText)
	}
}
