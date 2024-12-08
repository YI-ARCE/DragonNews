package sqlModel

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"unsafe"
	"yiarce/core/file"
	"yiarce/core/yorm"
)

var dbs *yorm.Db

var model string

var rootDir = ``

func LinkDb(config yorm.Config) {
	db, dbErr := yorm.ConnMysql(config)
	if dbErr != nil {
		fmt.Println(dbErr.Error())
		return
	}
	dbs = db
	r, err := db.QueryMap(`show tables`)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	s, _ := file.Get(`./module/sqlModel/table.model`)
	model = s.String()
	if rootDir == `` {
		f, _ := file.Get(`./go.mod`)
		rootDir = strings.Split(strings.Split(f.String(), "\n")[0], " ")[1]
		arr := strings.Split(rootDir, `\`)
		rootDir = arr[len(arr)-1]
	}
	s.Close()
	for i, m := range r {
		structure(m[`Tables_in_`+config.Database], i+1)
	}
}

func structure(name string, i int) {
	//tableName := NameToUp(name)
	r, err := dbs.QueryMap(`SELECT table_comment FROM information_schema.TABLES where TABLE_NAME = '` + name + `'`)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	remark := r[0][`TABLE_COMMENT`]
	r, err = dbs.QueryMap(`show FULL FIELDS FROM ` + name)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fStr := `// Structure ` + remark + "\n" + `type Structure struct{` + "\n"
	alias := ToAlias(name) + strconv.Itoa(i)
	constK := ``
	aliasK := ``
	for _, m := range r {
		tk := tableKey{}
		tk.name = m[`Field`]
		tk.isNull = m[`Null`] == `YES`
		tk.pri = m[`Key`] == `PRI`
		tkType, _ := checkType(m[`Type`])
		tk.types = tkType
		tk.remark = m[`Comment`]
		fStr += `	` + `// ` + tk.remark + "\n" + `	` + NameToUp(tk.name) + ` ` + tk.types + " `json:\"" + tk.name + `,omitempty"` + "`" + "\n"
		constK += `const ` + NameToUp(tk.name) + ` = ` + "`" + tk.name + "`" + ` // ` + tk.remark + "\n\n"
		aliasK += "// " + NameToUp(tk.name) + ` ` + tk.remark + "\n" + `func (a alias) ` + NameToUp(tk.name) + `() string {` + "\n	return *(*string)(&a) + " + NameToUp(tk.name) + "\n}\n\n"
	}
	fStr += `}`
	nStr := strings.ReplaceAll(model, `[dn:packageName]`, NameToUp(name, 1))
	nStr = strings.ReplaceAll(nStr, `[dn:rootDir]`, rootDir)
	nStr = strings.ReplaceAll(nStr, `[dn:alias]`, alias)
	nStr = strings.ReplaceAll(nStr, `[dn:tableName]`, name)
	nStr = strings.ReplaceAll(nStr, `[dn:constKey]`, constK[:len(constK)-2])
	nStr = strings.ReplaceAll(nStr, `[dn:Struct]`, fStr)
	nStr = strings.ReplaceAll(nStr, `[dn:aliasKey]`, aliasK[:len(aliasK)-2])
	file.Set(`table/`+NameToUp(name, 1), `table.go`, *(*[]byte)(unsafe.Pointer(&nStr)), os.O_CREATE|os.O_TRUNC, 0777)
}

func checkType(t string) (string, string) {
	ts := strings.Split(t, `(`)
	ts[1] = ts[1][:len(ts[1])-1]
	types := ``
	switch ts[0] {
	case `int`, `tinyint`:
		types = `int`
	case `char`, `varchar`, `text`:
		types = `string`
	case `decimal`:
		types = `float`
	default:
		types = `undefined`
	}
	return types, ts[1]
}

func NameToUp(name string, flag ...int) string {
	tName := strings.Split(name, `_`)
	str := ``
	if len(flag) > 0 {
		str = tName[0]
		tName = tName[1:]
	}
	for _, s := range tName {
		up := strings.ToUpper(s[0:1])
		str += up + s[1:]
	}
	return str
}

func ToAlias(name string) string {
	tName := strings.Split(name, `_`)
	str := ``
	for _, s := range tName {
		str += s[0:1]
	}
	return str
}
