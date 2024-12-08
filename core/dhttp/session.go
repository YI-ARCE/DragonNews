package dhttp

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

// session Session管理单元
var sessionStorage storage

func init() {
	sessionStorage = storage{pool: map[string]session{}}
}

func (sr *SessionReader) Get(key string) interface{} {
	cookie, err := sr.r.Cookie("dragon-cookie")
	if err != nil || cookie.Value == "" || sessionStorage.pool[cookie.Value].data[key] == nil {
		return nil
	}
	return sessionStorage.pool[cookie.Value].data[key]
}

func (sr *SessionReader) Set(key string, content interface{}) bool {
	cookie, err := sr.r.Cookie("dragon-cookie")
	if err != nil || cookie.Value == "" {
		cookies := http.Cookie{
			Name:     "dragon-cookie",
			Value:    sessionId(),
			Path:     "/",
			HttpOnly: true,
			MaxAge:   7200,
			Expires:  time.Now().Add(time.Duration(7200)),
		}
		http.SetCookie(*sr.w, &cookies)
		sessionStorage.pool[cookies.Value] = session{
			sessionId: cookies.Value,
			time:      time.Now().Unix(),
			maxTime:   7200,
			lock:      new(sync.Mutex),
			data:      map[string]interface{}{key: content},
		}
		sr.r.AddCookie(&cookies)
		return true
	}
	if sessionStorage.pool[cookie.Value].sessionId == "" {
		sessionStorage.pool[cookie.Value] = session{
			sessionId: cookie.Value,
			time:      time.Now().Unix(),
			maxTime:   7200,
			lock:      new(sync.Mutex),
			data:      map[string]interface{}{key: content},
		}
		return true
	}
	sessionStorage.pool[cookie.Value].lock.Lock()
	sessionStorage.pool[cookie.Value].data[key] = content
	defer sessionStorage.pool[cookie.Value].lock.Unlock()
	return true
}

func Clear() {

}

// 验证session是否过期,过期将清除
func verifySession() {

}

// 获取ID
func sessionId() string {
	bytes := make([]byte, 48)
	if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
		return ""
	}
	h := md5.New()
	h.Write([]byte(base64.URLEncoding.EncodeToString(bytes)))
	return hex.EncodeToString(h.Sum(nil))
}
