package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"strings"
)

func Link() {
	http.HandleFunc("/", socket)
	http.ListenAndServe(":55520", nil)
}

var conn *websocket.Conn

func Flag() bool {
	return conn == nil
}

func socket(w http.ResponseWriter, r *http.Request) {
	wu := websocket.Upgrader{}
	wu.ReadBufferSize = 1024
	wu.WriteBufferSize = 1024
	wu.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	c, err := wu.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	conn = c
	onConnect(r)
	defer c.Close()
f:
	for {
		t, p, _ := c.ReadMessage()
		switch t {
		case -1:
			break f
		case 2:
			break
		case 1:
			onMessage(string(p))
		default:
			fmt.Println("类型:", t, "消息:", p)
			break
		}
	}
}

func onConnect(r *http.Request) {
	fmt.Println(r.RemoteAddr[0:strings.IndexAny(r.RemoteAddr, ":")], "连接成功")
}

func onMessage(str string) {
	fmt.Println(str)
}

func Push(tag string, str string) {
	t := `{"type":"` + tag + `","data":` + str + `}`
	if conn != nil {
		conn.WriteMessage(2, []byte(t))
	}
}
