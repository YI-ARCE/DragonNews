package test

import (
	"testing"
	"yiarce/core/frame"
)

// TestErrorCreation 测试错误创建功能
func TestErrorCreation(t *testing.T) {
	// 测试创建HTTP错误
	httpErr := frame.NewError(frame.HttpError, "HTTP错误测试")
	if httpErr == nil {
		t.Error("NewError returned nil")
	}
	if !httpErr.IsApi {
		t.Error("HTTP error should have IsApi = true")
	}
	if httpErr.Message != "HTTP错误测试" {
		t.Errorf("Error message mismatch, expected: HTTP错误测试, got: %s", httpErr.Message)
	}

	// 测试创建框架错误
	frameErr := frame.NewError(frame.SelfError, "框架错误测试")
	if frameErr == nil {
		t.Error("NewError returned nil")
	}
	if !frameErr.IsFrame {
		t.Error("Frame error should have IsFrame = true")
	}
	if frameErr.Message != "框架错误测试" {
		t.Errorf("Error message mismatch, expected: 框架错误测试, got: %s", frameErr.Message)
	}
}

// TestErrorInterface 测试错误接口实现
func TestErrorInterface(t *testing.T) {
	// 创建一个错误
	err := frame.NewError(frame.HttpError, "错误测试")
	if err == nil {
		t.Error("NewError returned nil")
	}

	// 测试Error()方法
	errorMsg := err.Error()
	if errorMsg != "错误测试" {
		t.Errorf("Error() method returned wrong message, expected: 错误测试, got: %s", errorMsg)
	}
}
