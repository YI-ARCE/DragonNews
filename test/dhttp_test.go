package test

import (
	"testing"
	"yiarce/core/dhttp"
)

// TestServerCreation 测试服务器创建功能
func TestServerCreation(t *testing.T) {
	// 测试创建服务器
	server := dhttp.Server("0.0.0.0", 8080)
	if server == (dhttp.Server("", 0)) {
		t.Errorf("Server creation failed")
	}
}

// TestServerTLSCreation 测试TLS服务器创建功能
func TestServerTLSCreation(t *testing.T) {
	// 测试创建TLS服务器
	server := dhttp.ServerTLS("0.0.0.0", 8443, "cert.pem", "key.pem")
	if server == (dhttp.ServerTLS("", 0, "", "")) {
		t.Errorf("ServerTLS creation failed")
	}
}

// TestRestart 测试服务器重启功能
func TestRestart(t *testing.T) {
	// 测试重启服务器
	err := dhttp.Restart()
	if err != nil {
		t.Errorf("Restart failed: %v", err)
	}
}

// TestCreateToken 测试创建token功能
func TestCreateToken(t *testing.T) {
	// 测试正常创建token
	id := 1
	phone := "13800138000"
	token := dhttp.CreateToken(id, phone)
	if token == "" {
		t.Errorf("CreateToken failed")
	}

	// 测试无效参数
	emptyToken := dhttp.CreateToken(0, phone)
	if emptyToken != "" {
		t.Errorf("CreateToken should return empty string for invalid ID")
	}

	emptyToken = dhttp.CreateToken(id, "")
	if emptyToken != "" {
		t.Errorf("CreateToken should return empty string for empty phone")
	}
}

// TestDecryptToken 测试解密token功能
func TestDecryptToken(t *testing.T) {
	// 首先创建一个token
	id := 1
	phone := "13800138000"
	token := dhttp.CreateToken(id, phone)
	if token == "" {
		t.Errorf("CreateToken failed")
		return
	}

	// 测试解密token
	decryptedToken, err := dhttp.DecryptToken(token)
	if err != nil {
		t.Errorf("DecryptToken failed: %v", err)
		return
	}

	if decryptedToken.ID != id {
		t.Errorf("DecryptToken returned wrong ID, expected: %d, got: %d", id, decryptedToken.ID)
	}

	if decryptedToken.Phone != phone {
		t.Errorf("DecryptToken returned wrong phone, expected: %s, got: %s", phone, decryptedToken.Phone)
	}

	// 测试无效token
	_, err = dhttp.DecryptToken("")
	if err == nil {
		t.Errorf("DecryptToken should return error for empty token")
	}

	// 测试错误格式的token
	_, err = dhttp.DecryptToken("invalid_token")
	if err == nil {
		t.Errorf("DecryptToken should return error for invalid token format")
	}
}
