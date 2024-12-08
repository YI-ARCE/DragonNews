package tools

import (
	"dragonNews/module/system/windows"
	"fmt"
)

type funcName struct {
	name string
	f    func()
}

func Scanf(echo string, format string, i interface{}) {
	fmt.Print(echo + ": ")
	fmt.Scanf(format+"\n", i)
}

func Start() {
	windows.Clear()
	var arr []funcName
	arr = append(arr, funcName{
		name: "长度计算",
		f:    LenStr,
	})
	fmt.Println(`-----工具-----`)
	for i, f := range arr {
		fmt.Println(i+1, `.`, f.name)
	}
	num := 0
	Scanf("工具标号", "%d", &num)
	windows.Clear()
	arr[num-1].f()
	Start()
}

// LenStr 计算字符长度
func LenStr() {
	s := ``
	Scanf(`请输入需要计算的字符`, "%s", &s)
	fmt.Println(`长度:`, len(s))
	windows.WaitToArbitraryKey("按任意键返回")
}
