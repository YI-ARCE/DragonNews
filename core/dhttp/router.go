package dhttp

import (
	"fmt"
	"yiarce/core/frame"
)

var manager map[string]*routerConstruct

func init() {
	manager = map[string]*routerConstruct{}
}

func routerCreate(path string, f func(r *Dn), noAuth ...bool) {
	if len(noAuth) > 0 {
		manager[path] = &routerConstruct{f: f, auth: false}
	} else {
		manager[path] = &routerConstruct{f: f, auth: true}
	}
}

func Get(path string, f func(r *Dn), noAuth ...bool) {
	routerCreate(path+"_GET", f, noAuth...)
}

func Post(path string, f func(r *Dn), noAuth ...bool) {
	routerCreate(path+"_POST", f, noAuth...)
}

func Rule(path string, f func(r *Dn), noAuth ...bool) {
	routerCreate(path, f, noAuth...)
}

func execute(r *Dn) {
	defer func() {
		rs := recover()
		if rs != nil {
			if err, ok := rs.(frame.Error); ok {
				if len(err.ApiCourse) > 0 {
					for _, s := range err.ApiCourse {
						fmt.Println(s)
					}
				}
				if len(err.FrameCurse) > 0 {
					for _, s := range err.FrameCurse {
						fmt.Println(s)
					}
				}
				r.Write(200, `{"code":0,"msg":"`+err.Message+`"}`)
			} else {
				defer func() {
					rs = recover()
					err = rs.(frame.Error)
					if len(err.ApiCourse) > 0 {
						for _, s := range err.ApiCourse {
							fmt.Println(s)
						}
					}
					if len(err.FrameCurse) > 0 {
						for _, s := range err.FrameCurse {
							fmt.Println(s)
						}
					}
				}()
				frame.Errors(frame.HttpError, `http`, r, 2)
			}
		}
	}()
	var rf *routerConstruct
	if function, flag := manager[r.uri]; flag {
		rf = function
	} else if function, flag = manager[r.uri+"_"+r.method]; flag {
		rf = function
	} else {
		r.Write(200, `403`)
		return
	}
	if inject(r, rf.auth) {
		rf.f(r)
	}
}
