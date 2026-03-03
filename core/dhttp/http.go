package dhttp

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
	"yiarce/core/file"
	"yiarce/core/monitor"
	"yiarce/core/yorm"

	"gopkg.in/yaml.v2"
)

type server struct {
	host string
	port int
	cert string
	key  string
}

type CfgServer struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Config struct {
	Server CfgServer `yaml:"server"`
}

var dbs *yorm.Db

var injectFunc func(d *Dn, auth bool) bool

// Server 设置监听地址及端口
func Server(host string, port int) server {
	return server{host: host, port: port}
}

func ServerTLS(host string, port int, cert string, key string) server {
	return server{host: host, port: port, cert: cert, key: key}
}

// Listen 开启服务
func (s server) Listen() error {
	log.SetOutput(io.Discard)
	http.HandleFunc("/", parse)
	addr := s.host + ":" + strconv.FormatInt(int64(s.port), 10)
	if s.cert != `` && s.key != `` {
		err := http.ListenAndServeTLS(addr, s.cert, s.key, nil)
		if err != nil {
			return fmt.Errorf("HTTPS服务启动失败: %w", err)
		}
	} else {
		err := http.ListenAndServe(addr, nil)
		if err != nil {
			return fmt.Errorf("HTTP服务启动失败: %w", err)
		}
	}
	return nil
}

// Listen 用于默认配置启动服务
//
// 功能：使用默认配置启动HTTP服务，监听0.0.0.0:8080
func Listen() error {
	s := server{host: "0.0.0.0", port: 8080}
	return s.Listen()
}

func parse(w http.ResponseWriter, r *http.Request) {
	// 记录请求开始时间
	startTime := time.Now()

	defer func() {
		// 计算请求执行时间
		duration := time.Since(startTime).Milliseconds()
		// 记录请求信息
		monitor.RecordRequest(r.URL.Path, duration)

		if err := recover(); err != nil {
			// 记录错误
			fmt.Printf("[ERROR] HTTP请求处理异常: %v\n", err)
			// 增加错误计数
			monitor.RecordError()
			// 返回500错误
			w.Header().Set("Content-Type", "application/json;charset=utf-8")
			_, _ = w.Write([]byte(`{"code":500,"msg":"服务器内部错误","success":false}`))
		}
	}()

	// 检查请求是否为nil
	if r == nil {
		w.Header().Set("Content-Type", "application/json;charset=utf-8")
		_, _ = w.Write([]byte(`{"code":400,"msg":"无效的请求","success":false}`))
		return
	}

	if r.Method == "OPTIONS" {
		_, _ = w.Write([]byte{})
		return
	}

	newRequest(&w, r)
}

// Restart 重启监听
//
// 功能：重启HTTP服务，目前仅作为预留接口
func Restart() error {
	// 注意：HTTP服务重启需要特殊处理，这里仅作为预留接口
	// 实际重启可能需要停止当前服务并重新启动
	return nil
}

func Inject(f func(dn *Dn, auth bool) bool) {
	injectFunc = f
}

// CfgStart 根据启动服务
func CfgStart() {
	df, _ := file.Get(`./config/database.yaml`)
	dfc := yorm.Config{}
	err := yaml.Unmarshal(df.Byte(), &dfc)
	if err != nil {
		monitor.Debug(`yorm error`, err.Error())
		return
	}
	dbs, err = yorm.Connect(dfc)
	if err != nil {
		monitor.Debug(`yorm error`, err.Error())
		return
	}
	monitor.Debug(`yorm`, dfc.Type, `connect success`)
	f, _ := file.Get(`./config/frame.yaml`)
	c := Config{}
	err = yaml.Unmarshal(f.Byte(), &c)
	if err != nil {
		monitor.Debug(`http error`, err.Error())
		return
	}
	monitor.Debug(`http`, `listen ->`, c.Server.Host, `:`, strconv.Itoa(c.Server.Port))
	Server(c.Server.Host, c.Server.Port).Listen()
}
