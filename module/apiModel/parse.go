package apiModel

import (
	"os"
	"strings"
	"unsafe"
	"yiarce/core/file"
)

var p []apiStruct
var imports [][]string
var rootPackage string
var alias string
var packages []string

func Start() {
	defer func() {
		return
	}()
	f, _ := file.Get(`go.mod`)
	rootPackage = strings.Split(strings.Split(f.String(), "\n")[0], ` `)[1]
	if len(os.Args) > 1 {
		path := strings.ReplaceAll(os.Args[1], `api.go`, ``)
		path = path[:len(path)-1]
		dir := path[strings.Index(path, `app`):]
		arr := strings.Split(dir, `\`)
		packages = arr[1:]
		alias = arr[0] + toUpper2(packages)
		imports = append(imports, []string{alias, strings.ReplaceAll(rootPackage+`\`+dir, `\`, `/`)})
		openFile(path)
		if len(p) > 0 {
			arr := strings.Split(os.Args[1], `\`)
			out(path, arr[len(arr)-2])
		}
	} else {
		dir, _ := os.ReadDir(`./app`)
		for _, entry := range dir {
			openFile(`./app/` + entry.Name() + `/`)
			if len(p) > 0 {
				out(`./app/`+entry.Name(), entry.Name())
				p = []apiStruct{}
			}
		}
	}
	outRouter()
}

func outRouter() {
	str := "	// [dn:alias]\n"
	as := ``
	for _, a := range p {
		as = a.alias
		str += `	dhttp.` + toUpper(a.method) + "(`" + toUrlHead() + `/` + toLower(a.name[0]) + "`," + a.alias + `.` + a.name[0] + `)` + "\n"
	}
	str = strings.ReplaceAll(str, `[dn:alias]`, as)
	s := ``
	f, err := file.Get(`./router/api.go`)
	if err == nil {
		s = f.String()
		f.Close()
	}
	head := ``
	for _, i2 := range imports {
		if strings.Contains(s, i2[0]+` "`+i2[1]+`"`) {
			continue
		} else {
			head += `import ` + i2[0] + ` "` + i2[1] + `"` + "\n"
		}
	}
	index := strings.Index(s, "\n\n"+`func Register()`)
	if index > -1 {
		s = strings.ReplaceAll(s, "\n"+`func Register()`, head+"\n"+`func Register()`)
		index2 := strings.Index(s, "	// "+alias)
		if index2 > -1 {
			index3 := strings.Index(s[index2+2:], "//")
			if index3 > -1 {
				s = strings.ReplaceAll(s, s[index2:index3], str)
			} else {
				s = strings.ReplaceAll(s, s[index2:strings.Index(s, "\n}")], str)
			}
		} else {
			s = strings.ReplaceAll(s, "\n}", str+"\n}")
		}
	} else {
		s = `package router` + "\n\n" + `import "yiarce/core/dhttp"` + "\n" + head + "\n" + `func Register() {` + "\n	// ---- api Auto\n" + str + `}`
	}
	file.Set(`./router`, `register.go`, *(*[]byte)(unsafe.Pointer(&s)), os.O_TRUNC|os.O_CREATE, 0777)
}

func openFile(path string) {
	f, _ := file.Get(path + `/api.go`)
	str := f.String()
	for checkRemark(&str) {
	}
}

func checkRemark(str *string) bool {
	start := strings.Index(*str, `// `)
	if start < 0 {
		return false
	}
	end := strings.Index(*str, `func `)
	createStruct((*str)[start:end])
	end2 := strings.Index(*str, "\n}\n")
	*str = (*str)[end2+1:]
	return true
}

func createStruct(str string) {
	arr := strings.Split(str, "//\n")
	if len(arr) < 2 {
		return
	}
	arr1 := strings.Split(strings.ReplaceAll(arr[0], "\n", ``), ` `)
	api := apiStruct{}
	api.name = arr1[1:]
	start := strings.Index(arr[1], "-param 请求参数\n")
	if start < 0 {
		return
	}
	end := strings.Index(arr[1], "\n//	-method")
	if end < 0 {
		arr[1] += "\n" + `//	-method empty`
		end = strings.Index(arr[1], "\n//	-method")
	}
	for _, s := range strings.Split(arr[1][start+20:end], "\n") {
		sa := strings.Split(s, " ")
		if len(sa) < 4 {
			continue
		}
		api.column = append(api.column, sa[1:])
	}
	//fmt.Println(api.column)
	api.alias = alias
	api.method = strings.Split(strings.ReplaceAll(arr[1][end+4:], "\n", ``), ` `)[1]
	p = append(p, api)
}

func out(path string, packageName string) {
	str := `package ` + packageName + "\n\n"
	for _, a := range p {
		str += `// ` + a.name[0] + `Model ` + a.name[1] + "\n" + `type ` + a.name[0] + `Model` + ` struct {` + "\n"
		for _, s := range a.column {
			str += `	` + toUpper(s[0]) + ` ` + s[1] + " `json:\"" + s[0]
			if len(s) == 4 {
				str += `,omitempty`
			}
			str += `"` + "`" + `	// ` + s[2] + "\n"
		}
		str = str[:len(str)-1]
		str += "\n}\n\n"
	}
	file.Set(path, `model.go`, *(*[]byte)(unsafe.Pointer(&str)), os.O_CREATE|os.O_TRUNC, 0777)
}

func toUrlHead() string {
	str := ``
	for _, s := range packages {
		str += s + `/`
	}
	return str[:len(str)-1]
}

func toUpper2(arr []string) string {
	str := ``
	for _, s := range arr {
		if len(s) == 0 {
			continue
		}
		str += strings.ToUpper(s[0:1]) + s[1:]
	}
	return str
}

func toLower(str string) string {
	return strings.ToLower(str[0:1]) + str[1:]
}

func toUpper(str string) string {
	arr := strings.Split(str, `_`)
	str = ``
	for _, s := range arr {
		if len(s) == 0 {
			continue
		}
		str += strings.ToUpper(s[0:1]) + s[1:]
	}
	return str
}
