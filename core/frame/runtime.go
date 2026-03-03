package frame

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"yiarce/core"
	"yiarce/core/date"

	"github.com/mattn/go-colorable"
)

const (
	HttpError   = `http`
	SelfError   = `frame`
	UnknowError = `unknow`
)

type messageColor string

type PrintConfig string

const (
	PrintDisAbleDebugInfo = PrintConfig(`[dragon]PrintDisAbleDebugInfo`)
	DefaultColor          = messageColor("\033[36m")
	ErrorColor            = messageColor("\033[31m")
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
	fmt.Print("\033[1m")
}

func echoLog(packageName []string, path string, file string, line int) string {
	l := len(packageName)
	return `[ 📁 ` + path + `/` + file + `:` + strconv.Itoa(line) + ` ] ` + `🪧` + packageName[l-1] + `()`
}

// 报错分类
func sorts(err *Error, packageName string, path string, log string) {
	if len(path) > 3 && path[:3] == `app` {
		// 指向HTTP服务
		err.ApiCourse = append(err.ApiCourse, log)
	}
	if len(packageName) >= 11 && packageName[:11] == `yiarce/core` {
		err.FrameCourse = append(err.FrameCourse, log)
	}
	err.Course = append(err.Course, log)
}

// Errors 错误拦截处理
func Errors(tag string, msg string, d HttpF, index ...int) {
	err := Error{Tag: tag}
	if d != nil {
		err.RequestArgs = HttpRequest{
			Get:  d.GetAll(),
			Body: core.Replace2Empty(string(d.Body()), "\r\n", "\r", "\n"),
		}
	}
	i := 1
	if len(index) > 0 {
		i = index[0]
	}
	dl := len(directory)
	for {
		pc, codePath, codeLine, oks := runtime.Caller(i)
		if !oks {
			break
		}
		packageName := strings.Split(runtime.FuncForPC(pc).Name(), `.`)
		pathIndex := strings.LastIndex(codePath, `/`)
		dCodePath := ""
		if pathIndex > -1 {
			startIndex := strings.Index(codePath[:pathIndex], directory)
			if startIndex > -1 {
				dCodePath = codePath[startIndex+dl+1 : pathIndex]
			} else {
				i += 1
				continue
			}
		}
		fileName := ""
		if pathIndex > -1 && pathIndex < len(codePath)-1 {
			fileName = codePath[pathIndex+1:]
		}
		sorts(&err, packageName[0], dCodePath, echoLog(packageName, dCodePath, fileName, codeLine))
		if !strings.Contains(codePath, directory) {
			i += 1
			continue
		}
		i += 1
	}
	err.Message = msg
	panic(err)
}

func Prevent(tag string, msg string, index ...int) {
	if len(index) > 0 {
		Errors(tag, msg, nil, index[0])
	} else {
		Errors(tag, msg, nil, 2)
	}
}

// NewError 创建一个新的Error对象
func NewError(tag string, msg string) *Error {
	err := &Error{
		Tag:     tag,
		Message: msg,
	}
	return err
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
	color := DefaultColor
	var pt []interface{}
	for _, s := range i {
		switch s.(type) {
		case messageColor:
			color = s.(messageColor)
		case PrintConfig:
			switch s.(PrintConfig) {
			case PrintDisAbleDebugInfo:
				flag = false
			}
		default:
			pt = append(pt, s)
		}
	}
	if flag {
		pc, codePath, codeLine, _ := runtime.Caller(1)
		packageName := strings.Split(runtime.FuncForPC(pc).Name(), `.`)
		rootIndex := strings.LastIndex(codePath, directory)
		dCodePath := ``
		if rootIndex > -1 {
			dCodePath = codePath[rootIndex+len(directory)+1:]
		}
		str := echoPrintLocation(packageName, dCodePath, codeLine)
		out.Write([]byte(fmt.Sprintf("\033[1m\033[34m%s\033[33m%s\033[31m%s", parseDate()+`[ 🐉 DragonNews ]`, str[1], ` `+str[2]+"\n")))
	}
	out.Write([]byte(color))
	printParseData(pt...)
}

func parseDate() string {
	str := `[ `
	t := date.New()
	hStr := t.Hour()
	mStr := t.Minutes()
	h, _ := strconv.Atoi(hStr)
	m, _ := strconv.Atoi(mStr)
	h = h % 12 * 2
	if m > 29 {
		h += 1
	}
	str += dateTag[h] + ` ` + t.Custom(`Y-M-D H:I:S`) + ` ]`
	return str
}

func printParseData(i ...interface{}) {
	l := len(i)
	for il, i3 := range i {
		v, err := serialize(i3)
		if err != nil {
			out.Write([]byte(fmt.Sprintf("[ frame-error ] %s -> %v", err.Error(), i3)))
		} else {
			out.Write([]byte(v))
		}
		if il < l-1 {
			out.Write([]byte(` `))
		}
	}
	out.Write([]byte("\n"))
}
