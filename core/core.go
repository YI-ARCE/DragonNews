// Package core
// 此文件用于提供框架基础方法及配置
package core

import (
	"os"
	"strings"
)

// 框架默认根位置
var path string

var Database map[string]string

func init() {
	path, _ = os.Getwd()
}

// Path 框架默认根位置
func Path() string {
	return path
}

func Replace2Empty(str string, old ...string) string {
	for _, v := range old {
		str = strings.ReplaceAll(str, v, ``)
	}
	return str
}

func MCReplace2Empty(str *string, old ...string) {
	for _, v := range old {
		*str = strings.ReplaceAll(*str, v, ``)
	}
}
