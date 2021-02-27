package session

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"
	"net/http"
	"sync"
	"time"
)

type session struct {
	SessionId string
	Time      int64
	MaxTime   int
	Lock      *sync.Mutex
	Data      map[string]interface{}
}

type storage struct {
	Lock    sync.Mutex
	Session map[string]session
}

type Http struct {
	R *http.Request
	W *http.ResponseWriter
}

//接口
type Session interface {
	Set(key string, content interface{}) bool
	Get(key string) (interface{}, bool)
}

//manager Session管理单元
var manager storage

func init() {
	manager = storage{Session: map[string]session{}}
}

func (h *Http) Get(key string) (interface{}, bool) {
	cookie, err := h.R.Cookie("dragon-cookie")
	if err != nil || cookie.Value == "" {
		return nil, false
	}
	return manager.Session[cookie.Value].Data[key], true
}

func (h *Http) Set(key string, content interface{}) bool {
	cookie, err := h.R.Cookie("dragon-cookie")
	if err != nil || cookie.Value == "" {
		cookies := http.Cookie{
			Name:     "dragon-cookie",
			Value:    getId(),
			Path:     "/",
			HttpOnly: true,
			MaxAge:   7200,
			Expires:  time.Now().Add(time.Duration(7200)),
		}
		http.SetCookie(*h.W, &cookies)
		manager.Session[cookies.Value] = session{
			SessionId: cookies.Value,
			Time:      time.Now().Unix(),
			MaxTime:   7200,
			Lock:      new(sync.Mutex),
			Data:      map[string]interface{}{key: content},
		}
		h.R.AddCookie(&cookies)
		return true
	}
	if manager.Session[cookie.Value].SessionId == "" {
		manager.Session[cookie.Value] = session{
			SessionId: cookie.Value,
			Time:      time.Now().Unix(),
			MaxTime:   7200,
			Lock:      new(sync.Mutex),
			Data:      map[string]interface{}{key: content},
		}
		return true
	}
	manager.Session[cookie.Value].Lock.Lock()
	manager.Session[cookie.Value].Data[key] = content
	defer manager.Session[cookie.Value].Lock.Unlock()
	return true
}

func Clear() {

}

//验证session是否过期,过期将清除
func verifySession() {

}

//获取ID
func getId() string {
	bytes := make([]byte, 48)
	if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
		return ""
	}
	h := md5.New()
	h.Write([]byte(base64.URLEncoding.EncodeToString(bytes)))
	return hex.EncodeToString(h.Sum(nil))
}
