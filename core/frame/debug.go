package frame

import (
	"fmt"
	"github.com/mattn/go-colorable"
	"os"
	"runtime"
	"strconv"
	"strings"
	"yiarce/core/date"
)

const (
	HttpError = iota
	SelfError
	UnfriendError
	TokenError
	SignError
)

const (
	PrintDisAbleDebugInfo = `[dragon]PrintDisAbleDebugInfo`
)

var directory string

var dateTag = []string{`🕛`, `🕧`, `🕐`, `🕜`, `🕑`, `🕝`, `🕒`, `🕞`, `🕓`, `🕟`, `🕔`, `🕠`, `🕕`, `🕡`, `🕖`, `🕢`, `🕗`, `🕣`, `🕘`, `🕤`, `🕙`, `🕥`, `🕚`, `🕦`}

var out = colorable.NewColorableStdout()

func init() {
	path, _ := os.Getwd()
	directory = strings.ReplaceAll(path, `\`, `/`)
}

func SetPackageName(name string) {
	directory = name
}

func echoLog(packageName []string, path string, file string, line int) string {
	l := len(packageName)
	return `[📦 ` + packageName[0] + `][📁 ` + path + `][🪧 ` + packageName[l-1] + `()] ` + file + ` 第 ` + strconv.Itoa(line) + ` 行`
}

// 报错分类
func sorts(err *Error, packageName string, path string, log string) {
	if len(path) > 3 && path[:3] == `app` {
		// 指向HTTP服务
		err.ApiCourse = append(err.ApiCourse, log)
	}
	if len(packageName) >= 11 && packageName[:11] == `yiarce/core` {
		err.FrameCurse = append(err.FrameCurse, log)
	}
	err.Course = append(err.Course, log)
}

// Errors 错误拦截处理
func Errors(types int, msg string, h HttpF, index ...int) {
	err := Error{}
	i := 1
	if len(index) > 0 {
		i = index[0]
	}
	for {
		pc, codePath, codeLine, oks := runtime.Caller(i)
		if !oks {
			break
		}
		packageName := strings.Split(runtime.FuncForPC(pc).Name(), `.`)
		pathIndex := strings.LastIndex(codePath, `/`)
		dCodePath := strings.Replace(codePath[:pathIndex], directory, "", 1)
		if len(dCodePath) > 1 {
			dCodePath = dCodePath[1:]
		}
		sorts(&err, packageName[0], dCodePath, echoLog(packageName, dCodePath, codePath[pathIndex+1:], codeLine))
		if !strings.Contains(codePath, directory) {
			i += 1
			continue
		}
		i += 1
	}
	switch types {
	case HttpError:
		err.IsApi = true
		if h != nil {
			h.Write(200, `{"code":0,"msg":"`+msg+`"}`)
		}
	case SelfError:
		err.IsFrame = true
	default:
		err.IsFrame = false
		err.IsApi = false
	}
	err.Message = msg
	panic(err)
}

func Prevent(tag int, msg string, index ...int) {
	if len(index) > 0 {
		Errors(tag, msg, nil, index[0])
	} else {
		Errors(tag, msg, nil, 2)
	}
}

func EchoError(i interface{}) {
	err := i.(Error)
	Println(err.Message)
	for k, v := range err.Course {
		Println(PrintDisAbleDebugInfo, k, `:`, v)
	}
}

func echoPrintLocation(packageName []string, path string, line int) []string {
	l := len(packageName)
	//`[📁 ` + path + ` ][🪧 ` + packageName[l-1] + `()] ` + file + ` 第 ` + strconv.Itoa(line) + ` 行`
	return []string{
		`[ 📁 ` + path + ` ]`,
		`[ 🪧 ` + packageName[l-1] + ` ]`,
		path + `:` + strconv.Itoa(line),
	}
}

func Println(i ...interface{}) {
	flag := true
	if len(i) > 1 {
		c, ok := i[0].(string)
		if ok && c == PrintDisAbleDebugInfo {
			flag = false
			i = i[1:]
		}
	}
	pc, codePath, codeLine, _ := runtime.Caller(1)
	packageName := strings.Split(runtime.FuncForPC(pc).Name(), `.`)
	rootIndex := strings.LastIndex(codePath, directory)
	dCodePath := ``
	if rootIndex > -1 {
		dCodePath = codePath[rootIndex:]
	}
	str := echoPrintLocation(packageName, dCodePath, codeLine)
	if flag {
		out.Write([]byte(fmt.Sprintf("\033[1m\033[34m%s\033[33m%s\033[31m%s", parseDate()+`[ 🐉 DragonNews ]`, str[1], ` `+str[2]+"\n")))
	}
	out.Write([]byte("\033[36m "))
	printParseData(i...)
}

func parseDate() string {
	str := `[ `
	t := date.Date()
	h, _ := strconv.ParseInt(t.Hour(), 10, 64)
	m, _ := strconv.ParseInt(t.Minutes(), 10, 64)
	index := h
	if index > 0 {
		index -= 1
	}
	if m > 29 {
		index += 1
	}
	str += dateTag[index] + ` ` + date.Date().Custom(`Y-M-D H:I:S`) + ` ]`
	return str
}

func printParseData(i ...interface{}) {
	l := len(i)
	for il, i3 := range i {
		out.Write([]byte(serialize(i3)))
		if il < l {
			fmt.Print(` `)
		}
	}
	out.Write([]byte("\n"))
}
