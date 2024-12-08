package dhttp

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

type server struct {
	host string
	port int
	cert string
	key  string
}

// Server 设置监听地址及端口
func Server(host string, port int) server {
	return server{host: host, port: port}
}

func ServerTLS(host string, port int, cert string, key string) server {
	return server{host: host, port: port, cert: cert, key: key}
}

// Listen 开启服务
func (s server) Listen() {
	log.SetOutput(io.Discard)
	//FilePath, _ := os.Getwd()
	//http.Handle("/static/", http.FileServer(http.Dir(FilePath)))
	http.HandleFunc("/api/", parse)
	if s.cert != `` && s.key != `` {
		err := http.ListenAndServeTLS(s.host+":"+strconv.FormatInt(int64(s.port), 10), s.cert, s.key, nil)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	} else {
		//var err error
		//fmt.Println(`服务启动中...`)
		//go func() {
		//	time.Sleep(time.Second * 3)
		//	if err == nil {
		//		fmt.Println(`服务启动成功,打开浏览器`)
		//		e := exec.Command(`cmd`, `/c`, `start`, `http://127.0.0.1:55031/static/index.html`)
		//		rr := e.Run()
		//		if rr != nil {
		//			fmt.Println(rr.Error())
		//			os.Exit(0)
		//		}
		//	}
		//}()
		//err = http.ListenAndServe(s.host+":"+strconv.FormatInt(int64(s.port), 10), nil)
		//if err != nil {
		//	fmt.Println(err.Error())
		//	return
		//}
		http.ListenAndServe(s.host+":"+strconv.FormatInt(int64(s.port), 10), nil)
	}
}

// Listen 用于默认配置启动服务
func Listen() {

}

func parse(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		_, _ = w.Write([]byte{})
		return
	}
	newRequest(&w, r)
}

// Restart 重启监听
func Restart() {

}
