package curl

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	Json  = "application/json;charset=utf-8"
	Xml   = "application/xml"
	Plain = "text/plain"
	Html  = "text/html"
	Form  = "application/x-www-form-urlencoded"
	File  = "multipart/form-data;"
)

type Reply struct {
	Request    http.Request
	Status     string
	StatusCode int
	//请求头数据
	Header map[string][]string
	//解析类型为字符串时会存于此
	Input string
	//获取的Json数据,非纯字符json建议使用Body自行转换
	Data map[string]interface{}
	//源数据,Data或Input无法满足时可以使用此自行解析
	Body []byte
	//文件会存放于此
	File map[string][]byte
	//用于接收octet-stream方式提交的数据,其处理需自行实现逻辑
	Other interface{}
}

// Get url 地址
// timer共四个参数,默认都为5S,超出部分会忽略
// timer1表示整个请求周期时间
// timer2表示等待响应的时间
// timer3表示寻址超时时间
// timer4表示读写超时时间
func Get(url string, contentType string, header map[string]string, timer ...int64) (Reply, error) {
	// 验证URL参数
	if url == "" {
		return Reply{}, errors.New("url cannot be empty")
	}

	// 处理超时参数
	times := [4]int64{5, 5, 5, 5}
	for key, val := range timer {
		if key > 3 { // 索引从0开始，所以最多处理4个参数（索引0-3）
			break
		}
		times[key] = val
	}

	// 创建HTTP客户端
	curl := create(times)

	// 创建HTTP请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Reply{}, err
	}

	// 设置请求头
	if header != nil {
		for key, val := range header {
			req.Header.Set(key, val)
		}
	}

	// 设置Content-Type
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	// 发送请求
	result, err := curl.Do(req)
	if err != nil {
		return Reply{}, err
	}
	defer result.Body.Close()

	return resp(result), nil
}

// Post url 地址,
// contentType 编码类型,
// timer共四个参数,默认都为5S,超出部分会忽略,
// timer1表示整个请求周期时间,
// timer2表示等待响应的时间,
// timer3表示寻址超时时间,
// timer4表示读写超时时间
func Post(url string, contentType string, header map[string]string, body string, timer ...int64) (Reply, error) {
	// 验证URL参数
	if url == "" {
		return Reply{}, errors.New("url cannot be empty")
	}

	// 处理超时参数
	times := [4]int64{5, 5, 5, 5}
	for key, val := range timer {
		if key > 3 { // 索引从0开始，所以最多处理4个参数（索引0-3）
			break
		}
		times[key] = val
	}

	// 创建HTTP客户端
	curl := create(times)

	// 创建HTTP请求
	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		return Reply{}, err
	}

	// 设置请求头
	if header != nil {
		for key, val := range header {
			req.Header.Set(key, val)
		}
	}

	// 设置Content-Type
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	// 发送请求
	result, err := curl.Do(req)
	if err != nil {
		return Reply{}, err
	}
	defer result.Body.Close()

	return resp(result), nil
}

// 创建一个http客户端对象模拟请求
func create(times [4]int64) *http.Client {
	curl := &http.Client{
		//请求生命周期时间
		Timeout: time.Second * time.Duration(times[0]),
		Transport: &http.Transport{
			// 等待响应的时间
			ResponseHeaderTimeout: time.Second * time.Duration(times[1]),
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				//寻址实例,并设置寻址超时时间
				conn, err := net.DialTimeout(network, addr, time.Second*time.Duration(times[2]))
				if err != nil {
					return nil, err
				}
				//读写超时时间
				conn.SetDeadline(time.Now().Add(time.Second * time.Duration(times[3])))
				return conn, nil
			},
		},
	}
	return curl
}

func resp(r *http.Response) Reply {
	// 初始化响应对象
	Res := Reply{
		Header: make(map[string][]string),
		Data:   make(map[string]interface{}),
		File:   make(map[string][]byte),
	}

	// 检查响应对象是否为nil
	if r == nil {
		return Res
	}

	// 读取响应体
	var body []byte
	var err error
	if r.Body != nil {
		body, err = io.ReadAll(r.Body)
		if err != nil {
			// 记录错误，但继续处理其他字段
			Res.Body = []byte{}
		} else {
			Res.Body = body
		}
	}

	// 设置响应头
	if r.Header != nil {
		Res.Header = r.Header
	}

	// 解析响应内容
	if body != nil && len(body) > 0 {
		headerContentType := r.Header["Content-Type"]
		if len(headerContentType) > 0 {
			contentType := strings.Split(headerContentType[0], ";")
			if len(contentType) > 0 {
				switch contentType[0] {
				case "application/json":
					err := json.Unmarshal(body, &Res.Data)
					if err != nil {
						// JSON解析失败，设置为空map
						Res.Data = make(map[string]interface{})
					}
					break
				case "text/plain":
					Res.Input = string(body)
					break
				case "text/html":
					Res.Input = string(body)
					break
				case "application/xml":
					Res.Input = string(body)
					break
				case "application/octet-stream":
					Res.Other = body
					break
				default:
					break
				}
			}
		}
	}

	// 设置响应状态
	Res.Status = r.Status
	Res.StatusCode = r.StatusCode
	Res.Request = *r.Request

	return Res
}
