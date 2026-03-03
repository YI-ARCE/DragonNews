package log

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
	"yiarce/core"
	"yiarce/core/date"
)

// 日志级别定义
const (
	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"
	FatalLevel = "fatal"
)

// 默认日志级别
const DefaultLevel = InfoLevel

// 日志文件大小限制（100MB）
const MaxLogFileSize = 100 * 1024 * 1024

// 日志通道缓冲区大小
const LogChanBufferSize = 10000

type Log struct {
	host   string
	method string
	uri    string
	ip     string
	flag   bool
	msg    string
	end    string
	skip   int
	level  string
	// 上下文信息
	context map[string]string
}

// 日志异步写入通道
var logChan = make(chan string, LogChanBufferSize)

// 日志通道关闭标志
var logChanClosed = false

// 日志写入协程启动标志
var logWriterStarted = false

// 日志写入协程互斥锁
var logWriterMutex sync.Mutex

var rootPath = core.Path()

// 多次少于32字符时用此方法
//
//	normal用此方法即可
func strBuild(str string, s ...string) string {
	b := strings.Builder{}
	b.Write([]byte(str))
	for _, v := range s {
		b.Write([]byte(v))
	}
	str = b.String()
	return str
}

// 多次大于32字符拼接时用此方法
func strJoin(str string, s ...string) string {
	bs := []string{str}
	for _, v := range s {
		bs = append(bs, v)
	}
	return strings.Join(bs, "")
}

func getPath() string {
	flag := strings.Index(rootPath, "\\")
	if flag > -1 {
		rootPath = strings.Replace(rootPath, "\\", "/", -1) + "/"
	}
	return rootPath
}

// Init 初始化日志对象
//
// 参数：
//   - Host: 主机名
//   - Method: 请求方法
//   - Uri: 请求路径
//   - Ip: 客户端IP
//
// 返回值：
//   - *Log: 日志对象
func Init(Host string, Method string, Uri string, Ip string) *Log {
	l := Log{}
	l.host = Host
	l.method = Method
	l.uri = Uri
	l.ip = Ip
	l.skip = 1
	l.level = DefaultLevel
	l.context = make(map[string]string)
	return &l
}

// SetContext 设置上下文信息
//
// 参数：
//   - key: 上下文键
//   - value: 上下文值
//
// 返回值：
//   - *Log: 日志对象
func (l *Log) SetContext(key string, value string) *Log {
	if l.context == nil {
		l.context = make(map[string]string)
	}
	l.context[key] = value
	return l
}

// GetContext 获取上下文信息
//
// 参数：
//   - key: 上下文键
//
// 返回值：
//   - string: 上下文值
func (l *Log) GetContext(key string) string {
	if l.context == nil {
		return ""
	}
	return l.context[key]
}

// Debug 记录调试级别日志
//
// 参数：
//   - s: 日志内容
func (l *Log) Debug(s interface{}) {
	l.skip = 2
	l.Insert(s, "🔍debug")
}

// Info 记录信息级别日志
//
// 参数：
//   - s: 日志内容
func (l *Log) Info(s interface{}) {
	l.skip = 2
	l.Insert(s, "🔖info")
}

// Warn 记录警告级别日志
//
// 参数：
//   - s: 日志内容
func (l *Log) Warn(s interface{}) {
	l.skip = 2
	l.Insert(s, "⚠️warn")
}

// Insert 插入日志内容
//
// 参数：
//   - s: 日志内容
//   - tags: 日志标签
func (l *Log) Insert(s interface{}, tags ...string) {
	str, err := serialize(s)
	if err != nil {
		return
	}
	tag := "♾️default"
	if len(tags) > 0 {
		tag = tags[0]
	}
	_, name, line, _ := runtime.Caller(l.skip)
	name = strings.Replace(name, getPath(), "", 1)

	// 构建日志头
	if len(l.msg) < 1 {
		l.msg = strBuild("[", tag, "]", `[`+date.New().Custom("Y-M-D H:I:S"), `]`, "[🔔 ", l.method, "]-", "[📌 ", l.uri, "]-", "[💻 ", l.ip, "] -> ", name, ":", strconv.FormatInt(int64(line), 10), "\n", str, "\n")
	} else {
		l.msg = strJoin(l.msg, str, "\n")
	}

	// 添加上下文信息
	if len(l.context) > 0 {
		l.msg = strBuild(l.msg, "[Context]")
		for k, v := range l.context {
			l.msg = strBuild(l.msg, " ", k, "=", v)
		}
		l.msg = strBuild(l.msg, "\n")
	}

	l.skip = 1
}

// Error 记录错误级别日志
//
// 参数：
//   - s: 日志内容
func (l *Log) Error(s interface{}) {
	l.skip = 2
	l.Insert(s, "⚠️error")
}

// Success 记录成功级别日志
//
// 参数：
//   - s: 日志内容
func (l *Log) Success(s interface{}) {
	l.skip = 2
	l.Insert(s, "♻️success")
}

// Fatal 记录致命级别日志
//
// 参数：
//   - s: 日志内容
func (l *Log) Fatal(s interface{}) {
	l.skip = 2
	l.Insert(s, "💥fatal")
}

// build 构建日志对象
//
// 参数：
//   - tag: 日志标签
//   - str: 日志内容
//
// 返回值：
//   - *Log: 日志对象
func (l *Log) build(tag string, str string) *Log {
	_, name, line, _ := runtime.Caller(2)
	name = strings.Replace(name, getPath(), "", 1)

	// 构建日志内容
	l.msg = strBuild("[", name, "("+strconv.FormatInt(int64(line), 10), ")]", "[", tag, "]", "[", date.New().Custom("Y-M-D H:I:S"), "]\n", str, "\n")

	// 添加上下文信息
	if len(l.context) > 0 {
		ctxStr := strBuild("[Context]")
		for k, v := range l.context {
			ctxStr = strBuild(ctxStr, " ", k, "=", v)
		}
		ctxStr = strBuild(ctxStr, "\n")
		l.msg = strBuild(ctxStr, l.msg)
	}

	return l
}

// Debug 全局调试级别日志
//
// 参数：
//   - s: 日志内容
func Debug(s interface{}) {
	str, err := serialize(s)
	if err != nil {
		return
	}
	l := Log{}
	_ = l.build("🔍debug", str).Out()
}

// Info 全局信息级别日志
//
// 参数：
//   - s: 日志内容
func Info(s interface{}) {
	str, err := serialize(s)
	if err != nil {
		return
	}
	l := Log{}
	_ = l.build("🔖info", str).Out()
}

// Warn 全局警告级别日志
//
// 参数：
//   - s: 日志内容
func Warn(s interface{}) {
	str, err := serialize(s)
	if err != nil {
		return
	}
	l := Log{}
	l.level = WarnLevel
	_ = l.build("⚠️warn", str).Out()
}

// Error 全局错误级别日志
//
// 参数：
//   - s: 日志内容
func Error(s interface{}) {
	str, err := serialize(s)
	if err != nil {
		return
	}
	l := Log{}
	_ = l.build("⚠️error", str).Out()
}

// Fatal 全局致命级别日志
//
// 参数：
//   - s: 日志内容
func Fatal(s interface{}) {
	str, err := serialize(s)
	if err != nil {
		return
	}
	l := Log{}
	_ = l.build("💥fatal", str).Out()
}

// Default
//
//	普通标记输出
//	所有快捷输出无法保证写入顺序及其他一些不可预测的异常问题
//	参数 s
func Default(s interface{}) {
	str, err := serialize(s)
	if err != nil {
		return
	}
	l := Log{}
	_ = l.build("♾️default", str).Out()
}

// Success 成功的输出方法,仅限输出传入内容
//
// 参数：
//   - s: 日志内容
func Success(s interface{}) {
	str, err := serialize(s)
	if err != nil {
		return
	}
	l := Log{}
	_ = l.build("♻️success", str).Out()
}

// startLogWriter 启动日志写入协程
func startLogWriter() {
	logWriterMutex.Lock()
	defer logWriterMutex.Unlock()

	if logWriterStarted {
		return
	}

	logWriterStarted = true

	go func() {
		for !logChanClosed {
			select {
			case logMsg, ok := <-logChan:
				if !ok {
					// 通道关闭
					return
				}
				// 写入日志文件
				writeLogToFile(logMsg)
			case <-time.After(1 * time.Second):
				// 防止协程阻塞
			}
		}
	}()
}

// writeLogToFile 将日志写入文件
//
// 参数：
//   - logMsg: 日志内容
func writeLogToFile(logMsg string) {
	defer func() {
		if err := recover(); err != nil {
			// 记录错误，但不影响其他日志的写入
			fmt.Println("[ERROR] Log writeLogToFile异常:", err)
		}
	}()

	if len(logMsg) < 1 {
		return
	}

	path := rootPath + `/runtime`
	d := date.New()
	path = strBuild(path, "/log/", d.Custom("YM"), "/")

	// 创建日志目录
	err := os.MkdirAll(path, 0755)
	if err != nil {
		// 记录错误，但不影响其他日志的写入
		fmt.Println("Error creating log directory:", err.Error())
		return
	}

	filename := d.Custom("D") + ".txt"
	file, err := os.OpenFile(path+filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		// 记录错误，但不影响其他日志的写入
		fmt.Println("Error opening log file:", err.Error())
		return
	}

	// 确保文件关闭
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			// 记录错误，但不影响其他日志的写入
			fmt.Println("Error closing log file:", closeErr.Error())
		}
	}()

	// 写入日志内容
	_, err = file.WriteString(logMsg)
	if err != nil {
		// 记录错误，但不影响其他日志的写入
		fmt.Println("Error writing to log file:", err.Error())
		return
	}

	// 强制刷新缓冲区，确保日志写入磁盘
	if syncErr := file.Sync(); syncErr != nil {
		// 记录错误，但不影响其他日志的写入
		fmt.Println("Error syncing log file:", syncErr.Error())
	}

	// 检查日志文件大小，超过限制时进行轮转
	fileInfo, err := file.Stat()
	if err == nil && fileInfo.Size() > MaxLogFileSize {
		// 由于使用了defer，不需要手动关闭文件

		// 重命名为带时间戳的文件
		oldPath := path + filename
		newPath := path + filename + "." + d.Custom("His")
		err = os.Rename(oldPath, newPath)
		if err != nil {
			fmt.Println("Error rotating log file:", err.Error())
		}
	}
}

// Out
//
//	调用此方法将实时写入当前所有记录的日志
//	返回 error
func (l *Log) Out() error {
	if len(l.msg) < 1 {
		return nil
	}

	// 启动日志写入协程
	startLogWriter()
	logMsg := l.msg + l.end

	// 将日志发送到通道
	select {
	case logChan <- logMsg:

		// 日志已发送到通道
	default:
		// 通道已满，直接写入文件
		writeLogToFile(logMsg)
	}

	// 清空消息，帮助GC
	l.msg = ""
	l.end = ""
	return nil
}
