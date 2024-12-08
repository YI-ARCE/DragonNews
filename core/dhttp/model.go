package dhttp

import (
	"net/http"
	"sync"
	"yiarce/core/log"
	"yiarce/core/yorm"
)

type routerConstruct struct {
	f    func(d *Dn) // 执行的API
	auth bool        // 是否需要检查
}

type Dn struct {
	host     string               //请求的server域名
	header   map[string][]string  // 请求头数据
	method   string               // 请求类型
	ip       string               //客户端 IP
	uri      string               // uri 访问标识
	input    string               //文本格式数据
	body     []byte               // 原始的请求数据
	get      map[string]string    // get 提交的数据
	post     map[string]string    // post 提交的数据
	file     map[string][]byte    // 表单提交的文件会存放于此
	response *http.ResponseWriter // 响应
	request  *http.Request
	Log      *log.Log // 日志
	Token    Token
	*yorm.Db
}

type session struct {
	sessionId string
	time      int64
	maxTime   int
	lock      *sync.Mutex
	data      map[string]interface{}
}

type storage struct {
	lock sync.Mutex
	pool map[string]session
}

type SessionReader struct {
	r *http.Request
	w *http.ResponseWriter
}
