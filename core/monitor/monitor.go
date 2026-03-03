package monitor

import (
	"fmt"
	"runtime"
	"sync"
	"time"
	"yiarce/core/date"
	"yiarce/core/frame"
	"yiarce/core/log"
	"yiarce/core/timing"
)

// 系统监控数据结构
type SystemInfo struct {
	StartTime      string           // 系统启动时间
	Uptime         int64            // 系统运行时间（秒）
	GoVersion      string           // Go语言版本
	CPUCount       int              // CPU核心数
	MemoryUsage    uint64           // 内存使用量（字节）
	GoroutineCount int              // Goroutine数量
	RequestCount   int64            // 请求次数
	ErrorCount     int64            // 错误次数
	RequestTime    map[string]int64 // 请求时间统计
	Mutex          sync.RWMutex     // 互斥锁
}

// 全局监控实例
var systemInfo = SystemInfo{
	StartTime:      date.DateTime(),
	Uptime:         0,
	GoVersion:      runtime.Version(),
	CPUCount:       runtime.NumCPU(),
	MemoryUsage:    0,
	GoroutineCount: 0,
	RequestCount:   0,
	ErrorCount:     0,
	RequestTime:    make(map[string]int64),
	Mutex:          sync.RWMutex{},
}

var l *log.Log

// 初始化监控模块
func init() {
	// 启动监控定时器
	timing.Anonymous(func() bool {
		updateSystemInfo()
		return true
	}, time.Second*5).Start()
	l = log.Init(`-`, `-`, `frame`, `-`)
}

// 更新系统监控信息
func updateSystemInfo() {
	systemInfo.Mutex.Lock()
	defer systemInfo.Mutex.Unlock()

	// 更新运行时间
	t, _ := time.Parse("2006-01-02 15:04:05", systemInfo.StartTime)
	systemInfo.Uptime = time.Now().Unix() - t.Unix()

	// 更新内存使用情况
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	systemInfo.MemoryUsage = m.Alloc

	// 更新Goroutine数量
	systemInfo.GoroutineCount = runtime.NumGoroutine()
}

// RecordRequest 记录请求信息
func RecordRequest(path string, duration int64) {
	systemInfo.Mutex.Lock()
	defer systemInfo.Mutex.Unlock()

	systemInfo.RequestCount++
	systemInfo.RequestTime[path] += duration
}

// RecordError 记录错误信息
func RecordError() {
	systemInfo.Mutex.Lock()
	defer systemInfo.Mutex.Unlock()

	systemInfo.ErrorCount++
}

// GetSystemInfo 获取系统监控信息
func GetSystemInfo() SystemInfo {
	systemInfo.Mutex.RLock()
	defer systemInfo.Mutex.RUnlock()

	// 创建一个新的SystemInfo结构体，只复制需要的字段，不复制锁
	// 对RequestTime进行深拷贝，避免外部修改影响内部数据
	requestTimeCopy := make(map[string]int64)
	for path, duration := range systemInfo.RequestTime {
		requestTimeCopy[path] = duration
	}

	return SystemInfo{
		StartTime:      systemInfo.StartTime,
		Uptime:         systemInfo.Uptime,
		GoVersion:      systemInfo.GoVersion,
		CPUCount:       systemInfo.CPUCount,
		MemoryUsage:    systemInfo.MemoryUsage,
		GoroutineCount: systemInfo.GoroutineCount,
		RequestCount:   systemInfo.RequestCount,
		ErrorCount:     systemInfo.ErrorCount,
		RequestTime:    requestTimeCopy,
		// 不复制锁字段
		Mutex: sync.RWMutex{},
	}
}

// PrintSystemInfo 打印系统监控信息
func PrintSystemInfo() {
	info := GetSystemInfo()
	fmt.Printf("\n=== 系统监控信息 ===\n")
	fmt.Printf("启动时间: %s\n", info.StartTime)
	fmt.Printf("运行时间: %d 秒\n", info.Uptime)
	fmt.Printf("内存使用量: %d KB\n", info.MemoryUsage/1024)
	fmt.Printf("Goroutine数量: %d\n", info.GoroutineCount)
	fmt.Printf("请求次数: %d\n", info.RequestCount)
	fmt.Printf("错误次数: %d\n", info.ErrorCount)
	fmt.Printf("=================\n\n")
}

func Debug(tag string, s ...interface{}) {
	c := []interface{}{frame.PrintDisAbleDebugInfo, `[`, tag, `]`}
	for _, s2 := range s {
		c = append(c, s2)
	}
	frame.Println(c...)
}
