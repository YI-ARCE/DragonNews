package dhttp

import (
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"yiarce/core/log"
)

func newRequest(w *http.ResponseWriter, r *http.Request) {
	body := r.Body
	d := &Dn{}
	d.body, _ = io.ReadAll(body)
	d.header = r.Header
	d.get = make(map[string]string)
	d.post = make(map[string]string)
	for index, value := range r.URL.Query() {
		d.get[index] = value[0]
	}
	_ = r.ParseForm()
	if len(r.Header["Content-Type"]) > 0 {
		switch strings.Split(r.Header["Content-Type"][0], ";")[0] {
		case "application/x-www-form-urlencoded":
			for index, value := range r.PostForm {
				d.post[index] = value[0]
			}
			break
		case "multipart/form-data":
			_ = r.ParseMultipartForm(0)
			for index, value := range r.MultipartForm.Value {
				d.post[index] = value[0]
			}
			for index := range r.MultipartForm.File {
				_, handler, err := r.FormFile(index)
				if err != nil {
					continue
				}
				d.file = map[string][]byte{}
				d.file[index] = getFile(handler)
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
	uriIndex := strings.Index(r.RequestURI, "?")
	if uriIndex != -1 {
		d.uri = r.RequestURI[0:uriIndex]
	} else {
		d.uri = r.RequestURI
	}
	d.host = r.Host
	d.ip = r.RemoteAddr[0:strings.IndexAny(r.RemoteAddr, ":")]
	d.Log = log.Init(d.host, d.method, d.uri[1:], d.ip)
	d.response = w
	d.request = r
	execute(d)
}

// 读取文件的字节流
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
