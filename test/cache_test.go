package test

import (
	"testing"
	"time"
	cache "yiarce/core/cache/file"
)

// TestCacheSetGet 测试缓存的设置和获取功能
func TestCacheSetGet(t *testing.T) {
	// 测试数据
	key := "test_key"
	value := "test_value"
	expire := time.Now().Add(1 * time.Hour).Unix()

	// 设置缓存
	err := cache.Set(key, value, expire)
	if err != nil {
		t.Errorf("Set cache error: %v", err)
	}

	// 获取缓存
	result, err := cache.Get(key)
	if err != nil {
		t.Errorf("Get cache error: %v", err)
	}

	// 验证缓存值
	if result != value {
		t.Errorf("Get cache value error, expected: %v, got: %v", value, result)
	}
}

// TestCacheExpire 测试缓存的过期功能
func TestCacheExpire(t *testing.T) {
	// 测试数据
	key := "test_expire_key"
	value := "test_expire_value"
	expire := time.Now().Add(-1 * time.Hour).Unix() // 设置为已过期

	// 设置缓存
	err := cache.Set(key, value, expire)
	if err != nil {
		t.Errorf("Set cache error: %v", err)
	}

	// 获取缓存（应该返回错误，因为缓存已过期）
	result, err := cache.Get(key)
	if err == nil {
		t.Errorf("Get expired cache should return error, but got: %v", result)
	}
}

// TestCacheDelete 测试缓存的删除功能
func TestCacheDelete(t *testing.T) {
	// 测试数据
	key := "test_delete_key"
	value := "test_delete_value"
	expire := time.Now().Add(1 * time.Hour).Unix()

	// 设置缓存
	err := cache.Set(key, value, expire)
	if err != nil {
		t.Errorf("Set cache error: %v", err)
	}

	// 删除缓存
	err = cache.Delete(key)
	if err != nil {
		t.Errorf("Delete cache error: %v", err)
	}

	// 获取缓存（应该返回错误，因为缓存已删除）
	result, err := cache.Get(key)
	if err == nil {
		t.Errorf("Get deleted cache should return error, but got: %v", result)
	}
}

// TestCacheClear 测试缓存的清除功能
func TestCacheClear(t *testing.T) {
	// 测试数据
	key1 := "test_clear_key1"
	value1 := "test_clear_value1"
	key2 := "test_clear_key2"
	value2 := "test_clear_value2"
	expire := time.Now().Add(1 * time.Hour).Unix()

	// 设置缓存
	err := cache.Set(key1, value1, expire)
	if err != nil {
		t.Errorf("Set cache error: %v", err)
	}
	err = cache.Set(key2, value2, expire)
	if err != nil {
		t.Errorf("Set cache error: %v", err)
	}

	// 清除所有缓存
	err = cache.Clear()
	if err != nil {
		t.Errorf("Clear cache error: %v", err)
	}

	// 获取缓存（应该返回错误，因为缓存已清除）
	result, err := cache.Get(key1)
	if err == nil {
		t.Errorf("Get cleared cache should return error, but got: %v", result)
	}
	result, err = cache.Get(key2)
	if err == nil {
		t.Errorf("Get cleared cache should return error, but got: %v", result)
	}
}

// TestCacheCleanup 测试缓存的清理功能
func TestCacheCleanup(t *testing.T) {
	// 测试数据
	key := "test_cleanup_key"
	value := "test_cleanup_value"
	expire := time.Now().Add(1 * time.Hour).Unix()

	// 设置缓存
	err := cache.Set(key, value, expire)
	if err != nil {
		t.Errorf("Set cache error: %v", err)
	}

	// 清除所有缓存
	err = cache.Clear()
	if err != nil {
		t.Errorf("Clear cache error: %v", err)
	}

	// 验证缓存是否已清除
	result, err := cache.Get(key)
	if err == nil {
		t.Errorf("Get cache after clear should return error, but got: %v", result)
	}
}
