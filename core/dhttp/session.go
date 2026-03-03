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
	sessionStorage = storage{pool: map[string]session{}, lock: sync.RWMutex{}}
}

func (sr *SessionReader) Get(key string) interface{} {
	cookie, err := sr.r.Cookie("dragon-cookie")
	if err != nil || cookie.Value == "" {
		return nil
	}
	sessionStorage.lock.RLock()
	defer sessionStorage.lock.RUnlock()
	if session, ok := sessionStorage.pool[cookie.Value]; ok {
		if val, ok := session.data[key]; ok {
			return val
		}
	}
	return nil
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
			Expires:  time.Now().Add(time.Second * 7200),
		}
		http.SetCookie(*sr.w, &cookies)
		sessionStorage.lock.Lock()
		sessionStorage.pool[cookies.Value] = session{
			sessionId: cookies.Value,
			time:      time.Now().Unix(),
			maxTime:   7200,
			lock:      new(sync.Mutex),
			data:      map[string]interface{}{key: content},
		}
		sessionStorage.lock.Unlock()
		sr.r.AddCookie(&cookies)
		return true
	}
	sessionStorage.lock.Lock()
	if s, ok := sessionStorage.pool[cookie.Value]; !ok || s.sessionId == "" {
		sessionStorage.pool[cookie.Value] = session{
			sessionId: cookie.Value,
			time:      time.Now().Unix(),
			maxTime:   7200,
			lock:      new(sync.Mutex),
			data:      map[string]interface{}{key: content},
		}
		sessionStorage.lock.Unlock()
		return true
	}
	sessionStorage.lock.Unlock()

	// 更新现有session的数据
	sessionStorage.lock.RLock()
	s := sessionStorage.pool[cookie.Value]
	sessionStorage.lock.RUnlock()

	s.lock.Lock()
	s.data[key] = content
	s.lock.Unlock()
	return true
}

// Clear 清除所有会话
func Clear() {
	sessionStorage.lock.Lock()
	defer sessionStorage.lock.Unlock()
	// 清空会话池
	sessionStorage.pool = make(map[string]session)
}

// 验证session是否过期,过期将清除
func verifySession() {
	sessionStorage.lock.Lock()
	defer sessionStorage.lock.Unlock()

	// 获取当前时间戳
	now := time.Now().Unix()

	// 遍历所有会话，检查是否过期
	for sessionId, s := range sessionStorage.pool {
		// 检查会话是否过期
		if now-s.time > int64(s.maxTime) {
			// 会话已过期，删除会话
			delete(sessionStorage.pool, sessionId)
		}
	}
}

// StartSessionManager 启动会话管理器，定期清理过期会话
func StartSessionManager() {
	// 每5分钟清理一次过期会话
	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for range ticker.C {
			verifySession()
		}
	}()
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
