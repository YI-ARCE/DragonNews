package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"gopkg.in/yaml.v2"
)

// Manager 配置管理器
type Manager struct {
	configs      map[string]*Config
	configLock   sync.RWMutex
	cacheTTL     time.Duration
	maxCacheSize int
}

// Config 配置对象
type Config struct {
	Data      map[string]interface{}
	LoadedAt  time.Time
	ExpiresAt time.Time
}

// NewManager 创建新的配置管理器
//
// 参数：
//   - cacheTTL: 缓存过期时间
//   - maxCacheSize: 最大缓存大小（0表示无限制）
//
// 返回值：
//   - *Manager: 配置管理器实例
func NewManager(cacheTTL time.Duration, maxCacheSize int) *Manager {
	return &Manager{
		configs:      make(map[string]*Config),
		cacheTTL:     cacheTTL,
		maxCacheSize: maxCacheSize,
	}
}

// LoadConfig 从指定路径加载配置文件
//
// 参数：
//   - filePath: 配置文件路径
//
// 返回值：
//   - map[string]interface{}: 配置数据
//   - error: 错误信息
func (cm *Manager) LoadConfig(filePath string) (map[string]interface{}, error) {
	// 检查文件路径
	if filePath == "" {
		return nil, fmt.Errorf("config file path cannot be empty")
	}

	// 获取绝对路径
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	// 检查文件是否存在
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist: %s", absPath)
	}

	// 检查缓存
	if config, found := cm.getCachedConfig(absPath); found {
		return config.Data, nil
	}

	// 读取文件
	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// 检查文件内容
	if len(data) == 0 {
		return nil, fmt.Errorf("config file is empty")
	}

	// 解析YAML
	var rawConfig map[interface{}]interface{}
	if err := yaml.Unmarshal(data, &rawConfig); err != nil {
		return nil, fmt.Errorf("failed to parse yaml config: %w", err)
	}

	// 转换为map[string]interface{}
	configMap := convertMap(rawConfig)

	// 缓存配置
	expiresAt := time.Now().Add(cm.cacheTTL)
	cm.setCachedConfig(absPath, &Config{
		Data:      configMap,
		LoadedAt:  time.Now(),
		ExpiresAt: expiresAt,
	})

	return configMap, nil
}

// GetConfig 获取指定路径的配置
//
// 参数：
//   - filePath: 配置文件路径
//
// 返回值：
//   - map[string]interface{}: 配置数据
//   - error: 错误信息
func (cm *Manager) GetConfig(filePath string) (map[string]interface{}, error) {
	return cm.LoadConfig(filePath)
}

// GetNestedConfig 获取嵌套的配置值
//
// 参数：
//   - filePath: 配置文件路径
//   - keys: 嵌套键名
//
// 返回值：
//   - interface{}: 配置值
//   - error: 错误信息
func (cm *Manager) GetNestedConfig(filePath string, keys ...string) (interface{}, error) {
	config, err := cm.LoadConfig(filePath)
	if err != nil {
		return nil, err
	}

	// 检查是否提供了键名
	if len(keys) == 0 {
		return config, nil
	}

	current := config
	for i, key := range keys {
		// 检查键名是否为空
		if key == "" {
			return nil, fmt.Errorf("key at index %d cannot be empty", i)
		}

		if next, ok := current[key]; ok {
			if mapNext, isMap := next.(map[string]interface{}); isMap {
				current = mapNext
			} else {
				return next, nil
			}
		} else {
			return nil, fmt.Errorf("key not found: %s", key)
		}
	}

	return current, nil
}

// RefreshConfig 刷新指定路径的配置
//
// 参数：
//   - filePath: 配置文件路径
//
// 返回值：
//   - map[string]interface{}: 配置数据
//   - error: 错误信息
func (cm *Manager) RefreshConfig(filePath string) (map[string]interface{}, error) {
	// 获取绝对路径
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	cm.invalidateCache(absPath)
	return cm.LoadConfig(filePath)
}

// invalidateCache 使缓存失效
//
// 参数：
//   - filePath: 配置文件路径
func (cm *Manager) invalidateCache(filePath string) {
	cm.configLock.Lock()
	delete(cm.configs, filePath)
	cm.configLock.Unlock()
}

// getCachedConfig 获取缓存的配置
//
// 参数：
//   - filePath: 配置文件路径
//
// 返回值：
//   - *Config: 配置对象
//   - bool: 是否找到缓存
func (cm *Manager) getCachedConfig(filePath string) (*Config, bool) {
	cm.configLock.RLock()
	config, found := cm.configs[filePath]
	cm.configLock.RUnlock()

	if !found {
		return nil, false
	}

	// 检查是否过期
	if time.Now().After(config.ExpiresAt) {
		cm.invalidateCache(filePath)
		return nil, false
	}

	return config, true
}

// setCachedConfig 设置缓存的配置
//
// 参数：
//   - filePath: 配置文件路径
//   - config: 配置对象
func (cm *Manager) setCachedConfig(filePath string, config *Config) {
	cm.configLock.Lock()
	defer cm.configLock.Unlock()

	// 检查缓存大小
	if cm.maxCacheSize > 0 && len(cm.configs) >= cm.maxCacheSize {
		// 删除最早的缓存
		var oldestPath string
		var oldestTime time.Time

		for path, cfg := range cm.configs {
			if oldestPath == "" || cfg.LoadedAt.Before(oldestTime) {
				oldestPath = path
				oldestTime = cfg.LoadedAt
			}
		}

		if oldestPath != "" {
			delete(cm.configs, oldestPath)
		}
	}

	cm.configs[filePath] = config
}

// ClearCache 清除所有缓存
//
// 返回值：
//   - error: 错误信息
func (cm *Manager) ClearCache() error {
	cm.configLock.Lock()
	defer cm.configLock.Unlock()

	cm.configs = make(map[string]*Config)
	return nil
}

// GetCacheSize 获取当前缓存大小
//
// 返回值：
//   - int: 缓存大小
func (cm *Manager) GetCacheSize() int {
	cm.configLock.RLock()
	defer cm.configLock.RUnlock()

	return len(cm.configs)
}

// DefaultManager 默认配置管理器
var DefaultManager = NewManager(5*time.Minute, 100)

// Load 加载配置文件（使用默认管理器）
//
// 参数：
//   - filePath: 配置文件路径
//
// 返回值：
//   - map[string]interface{}: 配置数据
//   - error: 错误信息
func Load(filePath string) (map[string]interface{}, error) {
	return DefaultManager.LoadConfig(filePath)
}

// Get 获取配置（使用默认管理器）
//
// 参数：
//   - filePath: 配置文件路径
//
// 返回值：
//   - map[string]interface{}: 配置数据
//   - error: 错误信息
func Get(filePath string) (map[string]interface{}, error) {
	return DefaultManager.GetConfig(filePath)
}

// GetNested 获取嵌套配置（使用默认管理器）
//
// 参数：
//   - filePath: 配置文件路径
//   - keys: 嵌套键名
//
// 返回值：
//   - interface{}: 配置值
//   - error: 错误信息
func GetNested(filePath string, keys ...string) (interface{}, error) {
	return DefaultManager.GetNestedConfig(filePath, keys...)
}

// Refresh 刷新配置（使用默认管理器）
//
// 参数：
//   - filePath: 配置文件路径
//
// 返回值：
//   - map[string]interface{}: 配置数据
//   - error: 错误信息
func Refresh(filePath string) (map[string]interface{}, error) {
	return DefaultManager.RefreshConfig(filePath)
}

// ClearCache 清除所有缓存（使用默认管理器）
//
// 返回值：
//   - error: 错误信息
func ClearCache() error {
	return DefaultManager.ClearCache()
}

// GetCacheSize 获取当前缓存大小（使用默认管理器）
//
// 返回值：
//   - int: 缓存大小
func GetCacheSize() int {
	return DefaultManager.GetCacheSize()
}

// convertMap 将map[interface{}]interface{}转换为map[string]interface{}
func convertMap(input map[interface{}]interface{}) map[string]interface{} {
	output := make(map[string]interface{})
	for key, value := range input {
		if strKey, ok := key.(string); ok {
			switch v := value.(type) {
			case map[interface{}]interface{}:
				output[strKey] = convertMap(v)
			case []interface{}:
				output[strKey] = convertSlice(v)
			default:
				output[strKey] = v
			}
		}
	}
	return output
}

// convertSlice 将[]interface{}转换为[]interface{}，处理其中的map[interface{}]interface{}
func convertSlice(input []interface{}) []interface{} {
	output := make([]interface{}, len(input))
	for i, value := range input {
		switch v := value.(type) {
		case map[interface{}]interface{}:
			output[i] = convertMap(v)
		case []interface{}:
			output[i] = convertSlice(v)
		default:
			output[i] = v
		}
	}
	return output
}
