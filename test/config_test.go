package test

import (
	"os"
	"testing"
	"yiarce/core/config"
)

// TestConfigLoad 测试配置加载功能
func TestConfigLoad(t *testing.T) {
	// 创建临时配置文件
	configContent := `
# 数据库配置
mysql:
  host: "127.0.0.1"
  port: "3306"
  database: "test_db"
  username: "root"
  password: "root"
  character: "utf8mb4"
`

	configFile := "test_config.yaml"
	err := os.WriteFile(configFile, []byte(configContent), 0644)
	if err != nil {
		t.Errorf("Create config file error: %v", err)
		return
	}
	defer os.Remove(configFile)

	// 加载配置
	cfg, err := config.Load(configFile)
	if err != nil {
		t.Errorf("Load config error: %v", err)
		return
	}

	// 验证配置值
	if mysql, ok := cfg["mysql"].(map[string]interface{}); ok {
		if host, ok := mysql["host"].(string); !ok || host != "127.0.0.1" {
			t.Errorf("Invalid mysql.host, expected: 127.0.0.1, got: %v", host)
		}
		if port, ok := mysql["port"].(string); !ok || port != "3306" {
			t.Errorf("Invalid mysql.port, expected: 3306, got: %v", port)
		}
		if database, ok := mysql["database"].(string); !ok || database != "test_db" {
			t.Errorf("Invalid mysql.database, expected: test_db, got: %v", database)
		}
		if username, ok := mysql["username"].(string); !ok || username != "root" {
			t.Errorf("Invalid mysql.username, expected: root, got: %v", username)
		}
		if password, ok := mysql["password"].(string); !ok || password != "root" {
			t.Errorf("Invalid mysql.password, expected: root, got: %v", password)
		}
		if character, ok := mysql["character"].(string); !ok || character != "utf8mb4" {
			t.Errorf("Invalid mysql.character, expected: utf8mb4, got: %v", character)
		}
	} else {
		t.Error("Invalid mysql config structure")
	}
}

// TestConfigGetNested 测试获取嵌套配置功能
func TestConfigGetNested(t *testing.T) {
	// 创建临时配置文件
	configContent := `
server:
  host: "0.0.0.0"
  port: 8080
  ssl:
    enabled: false
    cert_file: "cert.pem"
    key_file: "key.pem"
`

	configFile := "test_config_nested.yaml"
	err := os.WriteFile(configFile, []byte(configContent), 0644)
	if err != nil {
		t.Errorf("Create config file error: %v", err)
		return
	}
	defer os.Remove(configFile)

	// 获取嵌套配置
	host, err := config.GetNested(configFile, "server", "host")
	if err != nil {
		t.Errorf("Get nested config error: %v", err)
		return
	}
	if host != "0.0.0.0" {
		t.Errorf("Invalid server.host, expected: 0.0.0.0, got: %v", host)
	}

	port, err := config.GetNested(configFile, "server", "port")
	if err != nil {
		t.Errorf("Get nested config error: %v", err)
		return
	}
	if port != 8080 {
		t.Errorf("Invalid server.port, expected: 8080, got: %v", port)
	}

	sslEnabled, err := config.GetNested(configFile, "server", "ssl", "enabled")
	if err != nil {
		t.Errorf("Get nested config error: %v", err)
		return
	}
	if sslEnabled != false {
		t.Errorf("Invalid server.ssl.enabled, expected: false, got: %v", sslEnabled)
	}
}

// TestConfigRefresh 测试刷新配置功能
func TestConfigRefresh(t *testing.T) {
	// 创建临时配置文件
	configContent := `
app:
  name: "test_app"
  version: "1.0.0"
`

	configFile := "test_config_refresh.yaml"
	err := os.WriteFile(configFile, []byte(configContent), 0644)
	if err != nil {
		t.Errorf("Create config file error: %v", err)
		return
	}
	defer os.Remove(configFile)

	// 第一次加载配置
	cfg1, err := config.Load(configFile)
	if err != nil {
		t.Errorf("Load config error: %v", err)
		return
	}

	// 修改配置文件
	updatedConfigContent := `
app:
  name: "test_app_updated"
  version: "1.0.1"
`
	err = os.WriteFile(configFile, []byte(updatedConfigContent), 0644)
	if err != nil {
		t.Errorf("Update config file error: %v", err)
		return
	}

	// 刷新配置
	cfg2, err := config.Refresh(configFile)
	if err != nil {
		t.Errorf("Refresh config error: %v", err)
		return
	}

	// 验证配置是否已更新
	if cfg1["app"].(map[string]interface{})["name"] == cfg2["app"].(map[string]interface{})["name"] {
		t.Error("Config should be updated after refresh")
	}
	if cfg2["app"].(map[string]interface{})["name"] != "test_app_updated" {
		t.Errorf("Invalid app.name after refresh, expected: test_app_updated, got: %v", cfg2["app"].(map[string]interface{})["name"])
	}
	if cfg2["app"].(map[string]interface{})["version"] != "1.0.1" {
		t.Errorf("Invalid app.version after refresh, expected: 1.0.1, got: %v", cfg2["app"].(map[string]interface{})["version"])
	}
}
