package reply

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"reflect"
	"strings"
	"yiarce/dragonnews/log"
	"yiarce/dragonnews/session"
)

type ContentType struct {
	Xml  string
	Json string
	Html string
	Text string
}

type Request struct {
	//请求的server域名
	Host string
	//请求头数据
	Header map[string][]string
	//请求类型
	Method string
	//客户端IP
	IP string
	//路由
	Uri string
	//文本格式数据
	Input string
	//获取的JSON数据,非纯字符信息建议使用BodyByte转换
	Body map[string]interface{}
	//自行处理的主体Body
	BodyByte []byte
	//get提交的数据
	Get map[string]string
	//post提交的数据
	Post map[string]string
	//表单提交的文件会存放于此
	File map[string][]byte
	//用于接收octet-stream方式提交的数据,其处理需自行实现逻辑
	Other interface{}
}

type Reply struct {
	Request Request
	W       http.ResponseWriter
	R       *http.Request
	Ct      ContentType
	Session session.Http
	Log     log.Log
}

//type reply interface {
//	SetHeader(key string,value string)
//	Return(status int,data interface{},tag ...string)
//}

func Start(w http.ResponseWriter, r *http.Request) Reply {
	body, _ := ioutil.ReadAll(r.Body)
	Req := Request{}
	Req.BodyByte = body
	Req.Header = r.Header
	headerContentType := r.Header["Content-Type"]
	Req.Get = make(map[string]string)
	Req.Post = make(map[string]string)
	for urlIndex, urlValue := range r.URL.Query() {
		Req.Get[urlIndex] = urlValue[0]
	}
	_ = r.ParseForm()
	if len(headerContentType) > 0 {
		contentType := strings.Split(headerContentType[0], ";")
		switch contentType[0] {
		case "application/x-www-form-urlencoded":
			for postIndex, postValue := range r.PostForm {
				Req.Post[postIndex] = postValue[0]
			}
			break
		case "multipart/form-data":
			_ = r.ParseMultipartForm(0)
			for formIndex, formValue := range r.MultipartForm.Value {
				Req.Post[formIndex] = formValue[0]
			}
			for fileIndex := range r.MultipartForm.File {
				_, handler, err := r.FormFile(fileIndex)
				if err != nil {
					continue
				}
				Req.File = map[string][]byte{}
				Req.File[fileIndex] = getFile(handler)
			}
			break
		case "application/json":
			err := json.Unmarshal(body, &Req.Body)
			if err != nil {
				Req.Body = nil
			}
			break
		case "text/plain":
			result := string(body[:])
			Req.Input = result
			break
		case "text/html":
			result := string(body[:])
			Req.Input = result
			break
		case "application/xml":
			result := string(body[:])
			Req.Input = result
			break
		case "application/octet-stream":
			Req.Other = body
			break
		default:
			break
		}
	}
	Req.Method = r.Method
	Req.Uri = r.RequestURI
	Req.Host = r.Host
	Req.IP = r.RemoteAddr[0:strings.IndexAny(r.RemoteAddr, ":")]
	return Reply{
		Request: Req,
		W:       w,
		R:       r,
		Ct: ContentType{
			Xml:  "application/xml; charset=utf-8",
			Json: "application/json; charset=utf-8",
			Html: "text/html; charset=utf-8",
			Text: "text/plain; charset=utf-8",
		},
		Session: session.Http{
			W: &w, R: r,
		},
		Log: log.Log{Host: Req.Host, Method: Req.Method, Uri: Req.Uri, IP: Req.IP},
	}
}

//读取文件的字节流
func getFile(handler *multipart.FileHeader) []byte {
	file, err := handler.Open()
	if err != nil {
		panic(err)
	}
	defer file.Close()
	bytes := make([]byte, handler.Size)
	_, _ = file.Read(bytes)
	return bytes
}

//设置响应头
func (reply Reply) SetHeader(key string, value string) {
	reply.W.Header().Set(key, value)
}

//返回数据,参数三为返回的格式,可以不填,默认为识别的对应内容
func (reply Reply) Return(status int, data interface{}, tag ...string) {
	var bytes []byte
	switch reflect.TypeOf(data).Kind() {
	case reflect.String:
		bytes = []byte(data.(string))
		reply.W.Header().Set("Content-Type", "text/plain; charset=utf-8")
	case reflect.Slice:
		bytes = data.([]byte)
		reply.W.Header().Set("Content-Type", "text/plain; charset=utf-8")
	case reflect.Array, reflect.Map, reflect.Struct:
		var err error
		bytes, err = json.Marshal(data)
		if err != nil {
			fmt.Print(err)
		}
		reply.W.Header().Set("Content-Type", "application/json; charset=utf-8")
	default:
		var err error
		bytes, err = json.Marshal(data)
		if err != nil {
			fmt.Print(err)
		}
		reply.W.Header().Set("Content-Type", "text/plain; charset=utf-8")
	}
	for _, value := range tag {
		reply.W.Header().Set("Content-Type", value)
	}
	reply.W.WriteHeader(status)
	_, err := reply.W.Write(bytes)
	if err != nil {
		panic(err)
	}
}
