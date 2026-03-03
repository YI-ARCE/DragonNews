// Package core
// 此文件用于提供框架基础方法及配置
package core

import (
	"os"
	"strings"
)

// 框架默认根位置
var path string

// Database 数据库配置信息
// 用于存储数据库连接相关的配置参数
var Database map[string]string

// init 初始化函数
//
// 功能：初始化框架根路径
func init() {
	path, _ = os.Getwd()
}

// Path 获取框架默认根位置
//
// 返回值：
//   - string: 框架根路径
//
// 功能：返回框架的根路径，用于定位配置文件、日志文件等资源
func Path() string {
	return path
}

// Replace2Empty 替换字符串中的指定内容为空字符串
//
// 参数：
//   - str: 原始字符串
//   - old: 要替换的字符串列表
//
// 返回值：
//   - string: 替换后的字符串
//
// 功能：将原始字符串中的指定内容替换为空字符串，返回替换后的结果
func Replace2Empty(str string, old ...string) string {
	for _, v := range old {
		str = strings.ReplaceAll(str, v, ``)
	}
	return str
}

// MCReplace2Empty 替换指针指向的字符串中的指定内容为空字符串
//
// 参数：
//   - str: 指向原始字符串的指针
//   - old: 要替换的字符串列表
//
// 功能：将指针指向的字符串中的指定内容替换为空字符串，直接修改原字符串
func MCReplace2Empty(str *string, old ...string) {
	for _, v := range old {
		*str = strings.ReplaceAll(*str, v, ``)
	}
}
