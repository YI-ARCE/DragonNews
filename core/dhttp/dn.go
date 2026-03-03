package dhttp

import (
	"encoding/json"
	"errors"
	"net/http"
	"yiarce/core/frame"
	_ "yiarce/core/yorm/mysql"
)

func (d *Dn) Session() SessionReader {
	return SessionReader{d.request, d.response}
}

// Host 访问地址
func (d *Dn) Host() string {
	return d.host
}

// Header 请求头
func (d *Dn) Header(key string) string {
	if values, ok := d.header[key]; ok && len(values) > 0 {
		return values[0]
	}
	return ""
}

// Headers 获取所有请求头内容
func (d *Dn) Headers() map[string][]string {
	return d.header
}

// Method 请求方法
func (d *Dn) Method() string {
	//我抄
	return d.method
}

// Ip 访问IP地址
func (d *Dn) Ip() string {
	return d.ip
}

// Uri 访问的路径
func (d *Dn) Uri() string {
	return d.uri
}

// Input 获取文本型数据,如在content-type为text时用此会获得原始的文本数据
func (d *Dn) Input() string {
	return d.input
}

// Body 原始的二进制数据
func (d *Dn) Body() []byte {
	return d.body
}

// File 获取文件流数据
func (d *Dn) File(key string) []byte {
	return d.file[key]
}

func (d *Dn) Get(key string) string {
	return d.get[key]
}

func (d *Dn) GetAll() map[string]string {
	return d.get
}

func (d *Dn) Post(key string) string {
	return d.post[key]
}

func (d *Dn) PostAll() map[string]string {
	return d.post
}

func (d *Dn) SetHeader(key string, value string) {
	(*d.response).Header().Set(key, value)
}

func (d *Dn) SetStatusCode(code int) {
	(*d.response).WriteHeader(code)
}

// Write 此方法将直接返回给客户端,无需接口return
func (d *Dn) Write(code int, data string, header ...map[string]string) {
	(*d.response).WriteHeader(code)
	if len(header) > 0 {
		for index, value := range header[0] {
			(*d.response).Header().Set(index, value)
		}
	}
	_, err := (*d.response).Write([]byte(data))
	if err != nil {
		d.Log.Error(err.Error())
		return
	}
}

// Json 应答
// 传入任何支持转换JSON的数据
//
//	例如 map array slice struct
//	非JSON格式将报出异常
func (d *Dn) Json(data interface{}, code ...int) error {
	var str []byte
	var err error

	// 检查response是否为nil
	if d.response == nil {
		return errors.New("response is nil")
	}

	// 尝试直接序列化数据
	if str, err = json.Marshal(data); err != nil {
		// 如果是字符串类型，直接使用
		if s, ok := data.(string); ok {
			str = []byte(s)
		} else {
			return errors.New("仅限可encode的数据和文本数据")
		}
	}

	// 设置Content-Type，添加charset=utf-8
	if (*d.response).Header().Get("Content-Type") == "" {
		(*d.response).Header().Set("Content-Type", "application/json;charset=utf-8")
	}

	// 设置状态码
	statusCode := http.StatusOK
	if len(code) > 0 {
		statusCode = code[0]
		// 检查状态码是否有效
		if statusCode < 100 || statusCode >= 600 {
			statusCode = http.StatusOK
		}
	}
	(*d.response).WriteHeader(statusCode)

	// 输出数据
	_, err = (*d.response).Write(str)
	if err != nil {
		d.Log.Error(err.Error())
		return err
	}
	return nil
}

// SuccessJson 成功响应
func (d *Dn) SuccessJson(data interface{}) error {
	// 构建成功响应格式
	successResp := map[string]interface{}{
		"code": 1,
		"msg":  "操作成功",
		"data": data,
	}
	return d.Json(successResp, 200)
}

// ErrorJson 错误响应
func (d *Dn) ErrorJson(code int, message string) error {
	// 构建错误响应格式
	errorResp := map[string]interface{}{
		"code":    code,
		"msg":     message,
		"success": false,
	}
	return d.Json(errorResp, code)
}

func (d *Dn) Data() map[string]string {
	var m map[string]string
	err := json.Unmarshal(d.body, &m)
	if err != nil {
		return make(map[string]string)
	}
	return m
}

func (d *Dn) BodyFormat(format interface{}, errMsg ...string) {
	err := json.Unmarshal(d.body, format)
	if err != nil {
		if len(errMsg) > 0 {
			frame.Errors(frame.HttpError, errMsg[0], d)
		} else {
			frame.Errors(frame.HttpError, `提交内容不正确,请检查`, d)
		}
	}
}

// OutByte 输出二进制数据
//
// 参数：
//   - data: 二进制数据
//   - contentType: 内容类型，默认为application/octet-stream
//   - code: HTTP状态码，默认为200
//
// 返回值：
//   - error: 错误信息
//
// 功能：输出二进制数据到客户端
func (d *Dn) OutByte(data []byte, contentType string, code ...int) error {
	// 验证参数
	if data == nil {
		return errors.New("data cannot be nil")
	}

	// 设置Content-Type
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	(*d.response).Header().Set("Content-Type", contentType)

	// 设置状态码
	if len(code) > 0 {
		(*d.response).WriteHeader(code[0])
	} else {
		(*d.response).WriteHeader(http.StatusOK)
	}

	// 输出数据
	_, err := (*d.response).Write(data)
	if err != nil {
		d.Log.Error(err.Error())
		return err
	}

	return nil
}

// OutFile 输出文件
//
// 参数：
//   - fileData: 文件数据
//   - filename: 文件名
//   - code: HTTP状态码，默认为200
//
// 返回值：
//   - error: 错误信息
//
// 功能：输出文件到客户端，设置Content-Disposition为attachment，提示用户下载
func (d *Dn) OutFile(fileData []byte, filename string, code ...int) error {
	// 验证参数
	if fileData == nil {
		return errors.New("file data cannot be nil")
	}

	// 设置Content-Type
	(*d.response).Header().Set("Content-Type", "application/octet-stream")
	// 设置Content-Disposition
	if filename != "" {
		(*d.response).Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")
	}

	// 设置状态码
	if len(code) > 0 {
		(*d.response).WriteHeader(code[0])
	} else {
		(*d.response).WriteHeader(http.StatusOK)
	}

	// 输出文件数据
	_, err := (*d.response).Write(fileData)
	if err != nil {
		d.Log.Error(err.Error())
		return err
	}

	return nil
}
