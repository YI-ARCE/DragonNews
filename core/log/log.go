package log

import (
	"errors"
	"os"
	"runtime"
	"strconv"
	"strings"
	"yiarce/core"
	"yiarce/core/date"
)

type Log struct {
	host   string
	method string
	uri    string
	ip     string
	flag   bool
	msg    string
	end    string
	skip   int
}

var rootPath = core.Path()

// 多次少于32字符时用此方法
//
//	正常用此方法即可
func strBuild(str string, s ...string) string {
	b := strings.Builder{}
	b.Write([]byte(str))
	for _, v := range s {
		b.Write([]byte(v))
	}
	str = b.String()
	return str
}

// 多次大于32字符拼接时用此方法
func strJoin(str string, s ...string) string {
	bs := []string{str}
	for _, v := range s {
		bs = append(bs, v)
	}
	return strings.Join(bs, "")
}

func getPath() string {
	flag := strings.Index(rootPath, "\\")
	if flag > -1 {
		rootPath = strings.Replace(rootPath, "\\", "/", -1) + "/"
	}
	return rootPath
}

func Init(Host string, Method string, Uri string, Ip string) *Log {
	l := Log{}
	l.host = Host
	l.method = Method
	l.uri = Uri
	l.ip = Ip
	l.skip = 1
	return &l
}

func (l *Log) Insert(s interface{}, tags ...string) {
	str, err := serialize(s)
	if err != nil {
		return
	}
	tag := "♾️default"
	if len(tags) > 0 {
		tag = tags[0]
	}
	_, name, line, _ := runtime.Caller(l.skip)
	name = strings.Replace(name, getPath(), "", 1)
	if len(l.msg) < 1 {
		l.msg = strBuild("---------- [💻 ", l.ip, "]-[📌 ", l.uri, "]-[🔔 ", l.method, "] ----------\n",
			"[", tag, "]", "[", name, "(", strconv.FormatInt(int64(line), 10), ")", "][",
			date.Date().Custom("Y-M-D H:I:S"), "]\n", str, "\n")
		//l.msg = "---------- [💻 " + l.ip + "]-[📌 " + l.uri + "]-[🔔 " + l.method + "] ----------\n" +
		//	"[" + tag + "]" + "[" + name + "(" + strconv.FormatInt(int64(line), 10) + ")" + "][" + common.date().Custom("Y-M-D H:I:S") + "]\n" + str + "\n"
	} else {
		l.msg = strBuild(l.msg, "[", tag, "]", "[", name, "("+strconv.FormatInt(int64(line), 10), ")", "][", date.Date().Custom("Y-M-D H:I:S"), "]\n", str, "\n")
		//l.msg += "[" + tag + "]" + "[" + name + "(" + strconv.FormatInt(int64(line), 10) + ")" + "][" + common.date().Custom("Y-M-D H:I:S") + "]\n" + str + "\n"
	}
	l.skip = 1
}

func (l *Log) Error(s interface{}) {
	l.skip = 2
	l.Insert(s, "⚠️error")
}

func (l *Log) Success(s interface{}) {
	l.skip = 2
	l.Insert(s, "♻️success")
}

func (l *Log) build(tag string, str string) *Log {
	_, name, line, _ := runtime.Caller(2)
	name = strings.Replace(name, getPath(), "", 1)
	l.msg = strBuild("[", name, "("+strconv.FormatInt(int64(line), 10), ")]", "[", tag, "][", date.Date().Custom("Y-M-D H:I:S"), "]\n", str, "\n")
	return l
}

// Default
//
//	普通标记输出
//	所有快捷输出无法保证写入顺序及其他一些不可预测的异常问题
//	参数 s
func Default(s interface{}) {
	str, err := serialize(s)
	if err != nil {
		return
	}
	l := Log{}
	_ = l.build("♾️default", str).Out()
}

func Error(s interface{}) {
	str, err := serialize(s)
	if err != nil {
		return
	}
	l := Log{}
	_ = l.build("⚠️error", str).Out()
}

// Success 成功的输出方法,仅限输出传入内容
func Success(s interface{}) {
	str, err := serialize(s)
	if err != nil {
		return
	}
	l := Log{}
	_ = l.build("♻️success", str).Out()
}

// Out
//
//	调用此方法将实时写入当前所有记录的日志
//	返回 error
func (l *Log) Out() error {
	if len(l.msg) < 1 {
		return nil
	}
	path := rootPath
	d := date.Date()
	path = strBuild(path, "/log/", d.Custom("YM"), "/")
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return errors.New("创建日志文件失败")
	}
	filename := d.Custom("D") + ".txt"
	file, err := os.OpenFile(path+filename, os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		return errors.New("创建日志文件失败")
	}
	_, err = file.WriteString(l.msg + l.end)
	if err != nil {
		return errors.New("日志文件数据写入失败")
	}
	err = file.Close()
	if err != nil {
		return err
	}
	return nil
}
