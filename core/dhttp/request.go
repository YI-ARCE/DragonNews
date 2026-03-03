package dhttp

import (
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"yiarce/core/date"
	"yiarce/core/log"
	"yiarce/core/monitor"
)

const (
	// MaxRequestBodySize 最大请求体大小（10MB）
	MaxRequestBodySize = 10 * 1024 * 1024
	// MaxFileSize 最大文件大小（5MB）
	MaxFileSize = 5 * 1024 * 1024
)

// newRequest 处理新的HTTP请求
//
// 参数：
//   - w: HTTP响应写入器
//   - r: HTTP请求
func newRequest(w *http.ResponseWriter, r *http.Request) {
	// 记录请求开始时间
	startTime := date.New()

	// 创建Dn对象
	d := &Dn{}
	defer func() {
		if err := recover(); err != nil {
			// 记录错误
			monitor.TagError(`error`, err)
			// 返回500错误
			(*w).Header().Set("Content-Type", "application/json;charset=utf-8")
			_, _ = (*w).Write([]byte(`{"code":500,"msg":"服务器内部错误","success":false}`))
		}
	}()

	// 限制请求体大小
	r.Body = http.MaxBytesReader(*w, r.Body, MaxRequestBodySize)

	// 读取请求体，捕获错误
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		if strings.Contains(err.Error(), "http: request body too large") {
			// 请求体过大
			(*w).Header().Set("Content-Type", "application/json;charset=utf-8")
			_, _ = (*w).Write([]byte(`{"code":413,"msg":"请求体过大","success":false}`))
			return
		}
		bodyBytes = []byte{}
	}
	d.body = bodyBytes
	d.header = r.Header

	// 预分配map容量，减少扩容开销
	d.get = make(map[string]string, len(r.URL.Query()))
	for index, value := range r.URL.Query() {
		if len(value) > 0 {
			d.get[index] = value[0]
		}
	}

	// 只有在需要时才解析表单
	if len(r.Header["Content-Type"]) > 0 {
		contentType := strings.Split(r.Header["Content-Type"][0], ";")[0]
		switch contentType {
		case "application/x-www-form-urlencoded":
			if err := r.ParseForm(); err == nil {
				d.post = make(map[string]string, len(r.PostForm))
				for index, value := range r.PostForm {
					if len(value) > 0 {
						d.post[index] = value[0]
					}
				}
			}
			break
		case "multipart/form-data":
			// 限制文件大小
			err := r.ParseMultipartForm(MaxFileSize)
			if err != nil {
				// 文件过大
				(*w).Header().Set("Content-Type", "application/json;charset=utf-8")
				_, _ = (*w).Write([]byte(`{"code":413,"msg":"文件过大","success":false}`))
				return
			}

			d.post = make(map[string]string, len(r.MultipartForm.Value))
			for index, value := range r.MultipartForm.Value {
				if len(value) > 0 {
					d.post[index] = value[0]
				}
			}

			if len(r.MultipartForm.File) > 0 {
				d.file = make(map[string][]byte, len(r.MultipartForm.File))
				for index := range r.MultipartForm.File {
					_, handler, err := r.FormFile(index)
					if err != nil {
						continue
					}

					// 检查单个文件大小
					if handler.Size > MaxFileSize {
						(*w).Header().Set("Content-Type", "application/json;charset=utf-8")
						_, _ = (*w).Write([]byte(`{"code":413,"msg":"文件过大","success":false}`))
						return
					}

					d.file[index] = getFile(handler)
				}
			}
			break
		case "text/plain", "text/html", "application/xml":
			result := string(d.body[:])
			d.input = result
			break
		default:
			break
		}
	}

	d.method = r.Method

	// 优化URI解析和安全检查
	uriIndex := strings.Index(r.RequestURI, "?")
	if uriIndex != -1 {
		d.uri = r.RequestURI[:uriIndex][1:]
	} else {
		d.uri = r.RequestURI[1:]
	}

	// 安全检查：防止路径遍历攻击
	if strings.Contains(d.uri, "..") {
		(*w).Header().Set("Content-Type", "application/json;charset=utf-8")
		_, _ = (*w).Write([]byte(`{"code":400,"msg":"无效的请求路径","success":false}`))
		return
	}

	d.host = r.Host

	// 优化IP地址提取，支持代理
	d.ip = getClientIP(r)

	// 延迟初始化日志，只有在需要时才创建
	d.Log = log.Init(d.host, d.method, d.uri, d.ip)
	d.response = w
	d.request = r
	monitor.Debug(`http`, `[`, startTime.HIS(`:`), `]`, `[ `+d.uri+` ]`, `[`, r.Method, `]`, `start`)

	// 执行请求处理
	execute(d)
	// 记录请求结束时间和处理耗时
	endTime := date.New()
	monitor.Debug(`http`, `[`, endTime.HIS(`:`), `]`, `[ `+d.uri+` ]`, `[ CLOSE ]`, `handle`, endTime.UnixMilli()-startTime.UnixMilli(), `ms`)
	d.Log.Out()
}

// getClientIP 获取客户端真实IP地址
//
// 参数：
//   - r: HTTP请求
//
// 返回值：
//   - string: 客户端IP地址
func getClientIP(r *http.Request) string {
	// 尝试从X-Forwarded-For头获取
	if xForwardedFor := r.Header.Get("X-Forwarded-For"); xForwardedFor != "" {
		// 多个代理时，第一个IP是客户端真实IP
		if parts := strings.Split(xForwardedFor, ","); len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}

	// 尝试从X-Real-IP头获取
	if xRealIP := r.Header.Get("X-Real-IP"); xRealIP != "" {
		return xRealIP
	}

	// 从RemoteAddr获取
	if colonIndex := strings.LastIndex(r.RemoteAddr, ":"); colonIndex > 0 {
		return r.RemoteAddr[:colonIndex]
	}

	return r.RemoteAddr
}

// getFile 读取文件的字节流
//
// 参数：
//   - handler: 文件头信息
//
// 返回值：
//   - []byte: 文件字节流
func getFile(handler *multipart.FileHeader) []byte {
	file, err := handler.Open()
	if err != nil {
		return []byte{}
	}
	defer file.Close()

	bytes := make([]byte, handler.Size)
	_, err = file.Read(bytes)
	if err != nil {
		return []byte{}
	}

	return bytes
}
