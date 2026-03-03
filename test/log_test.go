package test

import (
	"testing"
	"yiarce/core/log"
)

// TestLogInit 测试日志初始化功能
func TestLogInit(t *testing.T) {
	// 测试初始化日志对象
	l := log.Init("localhost", "GET", "test/path", "127.0.0.1")
	if l == nil {
		t.Error("Init returned nil")
	}
}

// TestLogSetLevel 测试设置日志级别功能
func TestLogSetLevel(t *testing.T) {
	// 初始化日志对象
	l := log.Init("localhost", "GET", "test/path", "127.0.0.1")
	if l == nil {
		t.Error("Init returned nil")
	}

	// 测试设置日志级别
	l = l.SetLevel(log.DebugLevel)
	if l == nil {
		t.Error("SetLevel returned nil")
	}

	l = l.SetLevel(log.InfoLevel)
	if l == nil {
		t.Error("SetLevel returned nil")
	}

	l = l.SetLevel(log.WarnLevel)
	if l == nil {
		t.Error("SetLevel returned nil")
	}

	l = l.SetLevel(log.ErrorLevel)
	if l == nil {
		t.Error("SetLevel returned nil")
	}

	l = l.SetLevel(log.FatalLevel)
	if l == nil {
		t.Error("SetLevel returned nil")
	}
}

// TestLogSetContext 测试设置上下文信息功能
func TestLogSetContext(t *testing.T) {
	// 初始化日志对象
	l := log.Init("localhost", "GET", "test/path", "127.0.0.1")
	if l == nil {
		t.Error("Init returned nil")
	}

	// 测试设置上下文信息
	l = l.SetContext("user_id", "123")
	if l == nil {
		t.Error("SetContext returned nil")
	}

	l = l.SetContext("request_id", "abc123")
	if l == nil {
		t.Error("SetContext returned nil")
	}
}

// TestLogMethods 测试日志记录方法
func TestLogMethods(t *testing.T) {
	// 初始化日志对象
	l := log.Init("localhost", "GET", "test/path", "127.0.0.1")
	if l == nil {
		t.Error("Init returned nil")
	}

	// 测试各种日志级别方法
	l.Debug("Debug log test")
	l.Info("Info log test")
	l.Warn("Warn log test")
	l.Error("Error log test")
	l.Fatal("Fatal log test")
	l.Success("Success log test")

	// 测试输出日志
	err := l.Out()
	if err != nil {
		t.Errorf("Out returned error: %v", err)
	}
}

// TestGlobalLogMethods 测试全局日志方法
func TestGlobalLogMethods(t *testing.T) {
	// 测试全局日志方法
	log.Debug("Global debug log test")
	log.Info("Global info log test")
	log.Warn("Global warn log test")
	log.Error("Global error log test")
	log.Fatal("Global fatal log test")
	log.Success("Global success log test")
	log.Default("Global default log test")
}

// TestLogLevelSetting 测试全局日志级别设置
func TestLogLevelSetting(t *testing.T) {
	// 测试设置全局日志级别
	log.SetGlobalLevel(log.DebugLevel)
	if log.GetGlobalLevel() != log.DebugLevel {
		t.Errorf("Global log level mismatch, expected: %s, got: %s", log.DebugLevel, log.GetGlobalLevel())
	}

	log.SetGlobalLevel(log.InfoLevel)
	if log.GetGlobalLevel() != log.InfoLevel {
		t.Errorf("Global log level mismatch, expected: %s, got: %s", log.InfoLevel, log.GetGlobalLevel())
	}

	log.SetGlobalLevel(log.WarnLevel)
	if log.GetGlobalLevel() != log.WarnLevel {
		t.Errorf("Global log level mismatch, expected: %s, got: %s", log.WarnLevel, log.GetGlobalLevel())
	}

	log.SetGlobalLevel(log.ErrorLevel)
	if log.GetGlobalLevel() != log.ErrorLevel {
		t.Errorf("Global log level mismatch, expected: %s, got: %s", log.ErrorLevel, log.GetGlobalLevel())
	}

	log.SetGlobalLevel(log.FatalLevel)
	if log.GetGlobalLevel() != log.FatalLevel {
		t.Errorf("Global log level mismatch, expected: %s, got: %s", log.FatalLevel, log.GetGlobalLevel())
	}
}
