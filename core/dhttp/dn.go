package dhttp

import (
	"encoding/json"
	"errors"
	"reflect"
	"unsafe"
	"yiarce/core/frame"
	"yiarce/core/yorm"
	_ "yiarce/core/yorm/mysql"
)

// 自定义注入,用于适应修改后的结构
//
//	auth 作为接口权限访问的检查,只有通过检查之后才会执行接口
//	若检查不通过,请在此执行响应结果并返回false
//	若不检查权限则直接返回true即可
func inject(d *Dn, auth bool) bool {
	if auth {
		// auth 检查不通过
		// d.Write(403, `{"msg":"访问接口没有权限","error":"no auth"}`)
		// return false
	}
	d.Token = DecryptToken(`M71q1pitQRsHd4Z13wYZEuenR8FSS53iBMd3XH+zy5igI2r11DzY8CyBZIF2nqNt`)
	// auth检查通过
	d.Db, _ = yorm.ConnMysql(yorm.Config{
		Host:     "127.0.0.1",
		Port:     "3306",
		Database: "sakura_disc",
		Username: "root",
		Password: "6353453wcR",
	})
	return true
}

func (d *Dn) Session() SessionReader {
	return SessionReader{d.request, d.response}
}

// Host 访问地址
func (d *Dn) Host() string {
	return d.host
}

// Header 请求头
func (d *Dn) Header(key string) string {
	return d.header[key][0]
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
	switch reflect.TypeOf(data).Kind() {
	case reflect.Slice, reflect.Array, reflect.Struct, reflect.Map:
		str, err = json.Marshal(data)
		if err != nil {
			return err
		}
		break
	case reflect.String:
		dStr := data.(string)
		str = *(*[]byte)(unsafe.Pointer(&dStr))
		break
	default:
		return errors.New("仅限可encode的数据和文本数据")
	}
	if (*d.response).Header().Get("Content-Type") == "" {
		(*d.response).Header().Set("Content-Type", "application/json")
	}
	if len(code) > 0 {
		(*d.response).WriteHeader(code[0])
	}
	_, err = (*d.response).Write(str)
	if err != nil {
		return err
	}
	return nil
}

func (d *Dn) Data() map[string]string {
	var m map[string]string
	_ = json.Unmarshal(d.body, &m)
	return m
}

func (d *Dn) BodyFormat(format interface{}, errMsg ...string) {
	err := json.Unmarshal(d.body, format)
	if err != nil {
		if len(errMsg) > 0 {
			frame.Errors(frame.HttpError, errMsg[0], d)
		} else {
			frame.Errors(frame.HttpError, `请求数据解析失败`, d)
		}
	}
}

func (d *Dn) OutByte() {

}

func (d *Dn) OutFile() {

}
