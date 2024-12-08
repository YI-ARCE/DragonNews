package yorm

import (
	"reflect"
	"strconv"
	"strings"
	"yiarce/core/frame"
)

type FormatStruct interface {
}

type Statement struct {
	Name   Name
	Alias  string
	Wheres []interface{}
	Fields string
	Orders []Orders
	Groups []string
	Joins  []Joins
	Pages  Pages
	Sql    bool
	Exec   bool
	d      *Db
}

type StatementData struct {
	s Statement
}

type Name struct {
	Name  string
	Alias string
}

type Wheres struct {
	Column  string
	Exp     string
	Content string
}

type Joins struct {
	JoinTable string
	JoinWay   string
	JoinWhere string
}

type Pages struct {
	Num  int
	Size int
}

type Orders struct {
	Column string
	Sort   string
}

var sample = Statement{
	Fields: "*",
	Pages:  Pages{0, 0},
}

func checkType(i interface{}) string {
	if i == nil {
		return ``
	}
	t := reflect.TypeOf(i)
	r := reflect.ValueOf(i)
	switch t.Kind() {
	case reflect.String:
		return `'` + r.String() + `'`
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(r.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(r.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(r.Float(), 'f', -1, 64)
	default:
		frame.Errors(frame.SelfError, "[yorm][sql-参数类型判断]Found build unsupported types : "+t.Kind().String(), nil)
	}
	return ""
}

func (d *StatementData) Table(name string, alias ...string) ModelTransfer {
	d.s.Name.Name = name
	l := len(alias)
	if l > 0 {
		d.s.Name.Alias = alias[0]
	}
	return d
}

func (d *StatementData) Where(column string, w ...interface{}) ModelTransfer {
	l := len(w)
	wh := Wheres{}
	switch l {
	case 1:
		wh.Column = column
		wh.Content = checkType(w[0])
		if wh.Content != `is null` && wh.Content != `is not null` {
			wh.Exp = `=`
		}
		if wh.Content != `` {
			d.s.Wheres = append(d.s.Wheres, wh)
		}
	case 2:
		wh.Column = column
		exp, flag := w[0].(string)
		if !flag {
			panic("[yorm][sql-where-func] w如果传递两个参数,第一个参数表示操作符,也就是说,只支持字符")
		}
		wh.Exp = exp
		if strings.Contains(strings.ToLower(wh.Exp), `in`) {
			wh.Content = `( ` + checkType(w[1]) + ` )`
		} else {
			wh.Content = checkType(w[1])
		}
		d.s.Wheres = append(d.s.Wheres, wh)
	default:
		if len(column) > 0 {
			d.s.Wheres = append(d.s.Wheres, column)
		}
	}
	return d
}

func (d *StatementData) Join(table string, condition string, link ...string) ModelTransfer {
	j := Joins{}
	j.JoinTable = table
	j.JoinWhere = condition
	if len(link) > 0 {
		j.JoinWay = link[0]
	}
	d.s.Joins = append(d.s.Joins, j)
	return d
}

func (d *StatementData) Field(column string) ModelTransfer {
	d.s.Fields = column
	return d
}

func (d *StatementData) Page(page int, size int) ModelTransfer {
	if page < 1 {
		panic("page can't be less than 1")
	}
	d.s.Pages.Num = page - 1
	d.s.Pages.Size = size
	return d
}

func (d *StatementData) Limit(size int) ModelTransfer {
	if size == 0 {
		panic("the size cannot be less than 1")
	}
	d.s.Pages.Num = 0
	d.s.Pages.Size = size
	return d
}

func (d *StatementData) Order(column string, sort string) ModelTransfer {
	o := Orders{
		column,
		sort,
	}
	d.s.Orders = append(d.s.Orders, o)
	return d
}

func (d *StatementData) Group(column string) ModelTransfer {
	d.s.Groups = append(d.s.Groups, column)
	return d
}

func (d *StatementData) FetchSql(flag ...bool) ModelTransfer {
	d.s.Sql = true
	if len(flag) > 0 && flag[0] {
		d.s.Exec = true
	}
	return d
}

// Find 转至驱动查询
func (d *StatementData) Find(fs ...FormatStruct) QueryResult {
	return d.s.d.driver.GetModel(d).Find(fs...)
}

// Select 转至驱动查询
func (d *StatementData) Select(fs ...FormatStruct) QueryResults {
	return d.s.d.driver.GetModel(d).Select(fs...)
}

// Update 转至驱动查询
func (d *StatementData) Update(i interface{}) ExecResult {
	return d.s.d.driver.GetModel(d).Update(i)
}

// Insert 转至驱动查询
func (d *StatementData) Insert(i interface{}) ExecResult {
	return d.s.d.driver.GetModel(d).Insert(i)
}

// Delete 转至驱动查询
func (d *StatementData) Delete() ExecResult {
	return d.s.d.driver.GetModel(d).Delete()
}

// Db 驱动需要此方法获取conn
func (d *StatementData) Db() *Db {
	return d.s.d
}

func (d *StatementData) GetStatement() *Statement {
	return &d.s
}
