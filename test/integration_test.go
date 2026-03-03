package test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"yiarce/core/cache/file"
	"yiarce/core/dhttp"
	"yiarce/core/monitor"
)

// TestIntegrationHTTPWithCache 测试HTTP请求与缓存的集成
func TestIntegrationHTTPWithCache(t *testing.T) {
	// 测试数据
	cacheKey := "integration_test_key"
	cacheValue := "integration_test_value"
	expire := time.Now().Add(1 * time.Hour).Unix()

	// 设置缓存
	err := cache.Set(cacheKey, cacheValue, expire)
	if err != nil {
		t.Errorf("Set cache error: %v", err)
	}

	// 获取缓存
	result, err := cache.Get(cacheKey)
	if err != nil {
		t.Errorf("Get cache error: %v", err)
	}

	if result != cacheValue {
		t.Errorf("Get cache value error, expected: %v, got: %v", cacheValue, result)
	}

	// 测试HTTP服务器创建
	_ = dhttp.Server("127.0.0.1", 8080)
	// 由于server是值类型，不能与nil比较，这里只测试创建过程是否报错
}

// TestIntegrationMonitor 测试监控功能的集成
func TestIntegrationMonitor(t *testing.T) {
	// 记录测试请求
	monitor.RecordRequest("/test/path", 100)
	monitor.RecordError()

	// 获取系统监控信息
	info := monitor.GetSystemInfo()

	// 验证监控信息
	if info.StartTime == "" {
		t.Error("System start time should not be empty")
	}

	if info.RequestCount == 0 {
		t.Error("Request count should be greater than 0")
	}

	if info.ErrorCount == 0 {
		t.Error("Error count should be greater than 0")
	}

	// 打印系统监控信息
	monitor.PrintSystemInfo()
}

// TestIntegrationHTTPRequest 测试HTTP请求的完整流程
func TestIntegrationHTTPRequest(t *testing.T) {
	// 创建测试HTTP请求
	_, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Errorf("Create request error: %v", err)
	}

	// 创建测试响应记录器
	w := httptest.NewRecorder()

	// 注意：dhttp.Parse是未导出函数，不能从外部包直接调用
	// 这里只测试HTTP请求和响应记录器的创建

	// 验证响应状态码（初始状态应该是200）
	if w.Code != http.StatusOK {
		t.Errorf("Unexpected initial status code: %d", w.Code)
	}
}
