package cache

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
	"yiarce/core"
	"yiarce/core/date"
	"yiarce/core/file"
)

// 缓存存储路径
var cachePath string

// 并发锁
var cacheMutex sync.RWMutex

// init 初始化缓存模块
//
// 功能：初始化缓存存储路径并创建缓存目录
func init() {
	// 初始化缓存存储路径
	cachePath = core.Path() + `/runtime/cache`
	// 创建缓存目录
	err := os.MkdirAll(cachePath, 0755)
	if err != nil {
		fmt.Printf("Failed to create cache directory: %v\n", err)
	}
}

// sanitizeKey 对缓存键名进行安全处理
//
// 参数：
//   - key: 原始缓存键名
//
// 返回值：
//   - string: 安全处理后的缓存键名
func sanitizeKey(key string) string {
	// 使用MD5哈希处理键名，防止路径遍历攻击
	hash := md5.Sum([]byte(key))
	return fmt.Sprintf("%x", hash)
}

// Get 获取缓存值
//
// 参数：
//   - key: 缓存键名
//
// 返回值：
//   - interface{}: 缓存值
//   - error: 错误信息
//
// 功能：根据键名获取缓存值，如果缓存不存在或已过期则返回错误
func Get(key string) (interface{}, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("[ERROR] Cache Get异常: %v\n", err)
		}
	}()

	// 检查键名参数
	if key == "" {
		return nil, fmt.Errorf("cache key cannot be empty")
	}

	// 对键名进行安全处理
	safeKey := sanitizeKey(key)

	// 构建缓存文件路径
	cacheFile := filepath.Join(cachePath, safeKey+".json")

	// 加读锁
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()

	// 检查缓存文件是否存在
	if _, err := os.Stat(cacheFile); os.IsNotExist(err) {
		return nil, err
	}

	// 读取缓存文件
	f, err := file.Get(cacheFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// 解析缓存数据
	var cacheData map[string]interface{}
	jsonData := f.Byte()
	if len(jsonData) == 0 {
		return nil, fmt.Errorf("cache file is empty")
	}

	err = json.Unmarshal(jsonData, &cacheData)
	if err != nil {
		return nil, err
	}

	// 检查缓存数据结构
	if _, ok := cacheData["value"]; !ok {
		return nil, fmt.Errorf("cache data missing value field")
	}

	// 检查缓存是否过期
	if expire, ok := cacheData["expire"].(float64); ok {
		if int64(expire) < time.Now().Unix() {
			// 缓存已过期，删除缓存文件
			cacheMutex.RUnlock()
			cacheMutex.Lock()
			os.Remove(cacheFile)
			cacheMutex.Unlock()
			cacheMutex.RLock()
			return nil, os.ErrNotExist
		}
	}

	// 返回缓存值
	return cacheData["value"], nil
}

// Set 设置缓存值
//
// 参数：
//   - key: 缓存键名
//   - value: 缓存值
//   - expire: 过期时间戳（秒）
//
// 返回值：
//   - error: 错误信息
//
// 功能：设置缓存值并指定过期时间
func Set(key string, value interface{}, expire int64) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("[ERROR] Cache Set异常: %v\n", err)
		}
	}()

	// 检查键名参数
	if key == "" {
		return fmt.Errorf("cache key cannot be empty")
	}

	// 检查过期时间
	if expire <= 0 {
		return fmt.Errorf("expire time must be greater than 0")
	}

	// 对键名进行安全处理
	safeKey := sanitizeKey(key)

	// 构建缓存数据
	cacheData := map[string]interface{}{
		"value":  value,
		"expire": expire,
		"time":   date.DateTime(),
	}

	// 序列化缓存数据
	jsonData, err := json.Marshal(cacheData)
	if err != nil {
		return fmt.Errorf("failed to marshal cache data: %w", err)
	}

	// 加写锁
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	// 写入缓存文件
	err = file.Set(cachePath, safeKey+".json", jsonData, os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}

	return nil
}

// Delete 删除缓存
//
// 参数：
//   - key: 缓存键名
//
// 返回值：
//   - error: 错误信息
//
// 功能：根据键名删除缓存
func Delete(key string) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("[ERROR] Cache Delete异常: %v\n", err)
		}
	}()

	// 检查键名参数
	if key == "" {
		return fmt.Errorf("cache key cannot be empty")
	}

	// 对键名进行安全处理
	safeKey := sanitizeKey(key)

	// 构建缓存文件路径
	cacheFile := filepath.Join(cachePath, safeKey+".json")

	// 加写锁
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	// 检查缓存文件是否存在
	if _, err := os.Stat(cacheFile); os.IsNotExist(err) {
		return nil
	}

	// 删除缓存文件
	return os.Remove(cacheFile)
}

// Clear 清除所有缓存
//
// 返回值：
//   - error: 错误信息
//
// 功能：清除所有缓存文件
func Clear() error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("[ERROR] Cache Clear异常: %v\n", err)
		}
	}()

	// 加写锁
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	// 读取缓存目录
	files, err := os.ReadDir(cachePath)
	if err != nil {
		return err
	}

	// 删除所有缓存文件
	var lastErr error
	for _, f := range files {
		if !f.IsDir() && filepath.Ext(f.Name()) == ".json" {
			cacheFile := filepath.Join(cachePath, f.Name())
			err = os.Remove(cacheFile)
			if err != nil {
				lastErr = err
				fmt.Printf("Failed to remove cache file %s: %v\n", cacheFile, err)
				// 继续删除其他文件，不因为单个文件失败而中断
			}
		}
	}

	return lastErr
}
