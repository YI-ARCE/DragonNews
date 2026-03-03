package dhttp

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"yiarce/core/frame"
	"yiarce/core/monitor"
)

var (
	manager     map[string]*routerConstruct
	routerMutex sync.RWMutex
)

func init() {
	manager = map[string]*routerConstruct{}
}

// routerCreate 创建路由
//
// 参数：
//   - path: 路由路径
//   - f: 路由处理函数
//   - auth: 是否需要权限检查
func routerCreate(types string, f func(r *Dn), auth ...bool) {
	if f == nil {
		fmt.Println("[ERROR] 路由处理函数不能为空")
		return
	}

	checkAuth := true
	if len(auth) > 0 {
		checkAuth = auth[0]
	}

	routerMutex.Lock()
	defer routerMutex.Unlock()
	path := getFuncPkg(f)
	monitor.Debug(`router`, `register api -> [ `+types+` ]`, path, `|`, `auth =`, checkAuth)
	path = path + `_` + types
	manager[path] = &routerConstruct{f: f, auth: checkAuth}
}

// Get 注册GET请求路由
//
// 参数：
//   - path: 路由路径
//   - f: 路由处理函数
//   - auth: 是否需要权限检查，默认需要权限检查
func Get(f func(r *Dn), auth ...bool) {

	routerCreate(`GET`, f, auth...)
}

// Post 注册POST请求路由
//
// 参数：
//   - path: 路由路径
//   - f: 路由处理函数
//   - auth: 是否需要权限检查，默认需要权限检查
func Post(f func(r *Dn), auth ...bool) {
	routerCreate(`POST`, f, auth...)
}

// Rule 注册无方法区分的路由
//
// 参数：
//   - path: 路由路径
//   - f: 路由处理函数
//   - noAuth: 是否需要权限检查
func Rule(f func(r *Dn), auth ...bool) {
	routerCreate(``, f, auth...)
}

// execute 执行路由处理函数
//
// 参数：
//   - r: Dn对象
func execute(r *Dn) {
	defer func() {
		if rs := recover(); rs != nil {
			if frameErr, ok := rs.(frame.Error); ok {
				// 记录错误信息
				r.Log.Error(frameErr)
				// 使用统一的错误响应格式
				errorResp := map[string]interface{}{
					"code":    0,
					"msg":     frameErr.Message,
					"success": false,
				}
				monitor.Error(frameErr)
				r.Json(errorResp)
			} else {
				// 处理非预期的错误
				defer func() {
					if rs := recover(); rs != nil {
						if frameErr, ok := rs.(frame.Error); ok {
							// 记录错误信息
							r.Log.Error(frameErr)
							monitor.Error(frameErr)
						}
					}
				}()
				// 使用统一的错误响应格式
				errorResp := map[string]interface{}{
					"code":    500,
					"msg":     "服务器内部错误",
					"success": false,
				}
				r.Json(errorResp, 500)
				frame.Errors(frame.UnknowError, fmt.Sprintf("%v", rs), r, 3)
			}
		}
	}()

	// 检查请求对象是否为nil
	if r == nil {
		return
	}

	// 优化路由查找逻辑
	var rf *routerConstruct
	var ok bool

	// 加读锁
	routerMutex.RLock()
	// 首先查找精确匹配的路由
	if rf, ok = manager[r.uri]; ok {
		// 路由存在
	} else if rf, ok = manager[r.uri+"_"+r.method]; ok {
		// 查找带HTTP方法后缀的路由
	}

	// 释放读锁
	routerMutex.RUnlock()

	if !ok {
		// 路由不存在，返回404
		notFoundResp := map[string]interface{}{
			"code":    404,
			"msg":     "路由不存在",
			"success": false,
		}
		r.Json(notFoundResp)
		return
	}

	// 检查路由处理函数是否为nil
	if rf.f == nil {
		errorResp := map[string]interface{}{
			"code":    500,
			"msg":     "路由处理函数未定义",
			"success": false,
		}
		r.Json(errorResp, 500)
		return
	}

	if injectFunc != nil && injectFunc(r, rf.auth) {
		r.Db = dbs
		rf.f(r)
	}
}

func getFuncPkg(fn interface{}) string {
	addr := reflect.ValueOf(fn).Pointer()
	f := runtime.FuncForPC(addr)
	arr := strings.Split(f.Name(), `.`)
	funcName := arr[1]
	arr = strings.Split(arr[0], `/`)
	if arr[1] != `app` {
		panic(`请使用app下的包函数注册为handle！`)
	}
	path := ``
	arr = arr[2:]
	for _, s := range arr {
		path = path + `/` + s
	}
	return path[1:] + `/` + strings.ToLower(funcName[0:1]) + funcName[1:]
}
