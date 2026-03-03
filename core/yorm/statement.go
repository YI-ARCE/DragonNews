package yorm

import (
	"reflect"
	"strconv"
	"strings"
	"yiarce/core/frame"
)

// FormatStruct 格式化结构接口
type FormatStruct interface {
}

// Statement SQL语句结构
type Statement struct {
	Name   Name          // 表名
	Alias  string        // 别名
	Wheres []interface{} // 查询条件
	Fields string        // 查询字段
	Orders []Orders      // 排序
	Groups []string      // 分组
	Joins  []Joins       // 连接
	Pages  Pages         // 分页
	SQL    bool          // 是否返回SQL
	Exec   bool          // 是否执行
	d      *Db           // 数据库连接
}

// StatementData 语句数据结构
type StatementData struct {
	s Statement
}

func (d *StatementData) Sum(column string) SumResult {
	return d.s.d.driver.GetModel(d).Sum(column)
}

// Name 表名结构
type Name struct {
	Name  string // 表名
	Alias string // 别名
}

// Wheres 查询条件结构
type Wheres struct {
	Column  string // 列名
	Exp     string // 操作符
	Content string // 内容
}

// Joins 连接结构
type Joins struct {
	JoinTable string // 连接表
	JoinWay   string // 连接方式
	JoinWhere string // 连接条件
}

// Pages 分页结构
type Pages struct {
	Num  int // 偏移量
	Size int // 大小
}

// Orders 排序结构
type Orders struct {
	Column string // 排序列
	Sort   string // 排序方式
}

// sample 示例语句
var sample = Statement{
	Fields: "*",
	Pages:  Pages{0, 0},
}

// checkType 检查参数类型并转换为字符串
//
// 参数：
//   - i: 任意类型的参数
//
// 返回值：
//   - string: 转换后的字符串
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
		frame.Errors(_tag, "类型转换失败: "+t.Kind().String(), nil)
	}
	return ""
}

// Table 设置表名
//
// 参数：
//   - name: 表名
//   - alias: 表别名
//
// 返回值：
//   - ModelTransfer: 模型转换接口
func (d *StatementData) Table(name string, alias ...string) ModelTransfer {
	// 检查表名参数
	if name == "" {
		frame.Errors(_tag, "查询表明不能为空", nil)
		return d
	}

	d.s.Name.Name = name
	l := len(alias)
	if l > 0 {
		// 检查别名参数
		if alias[0] == "" {
			frame.Errors(_tag, "查询表明不能为空", nil)
			return d
		}
		d.s.Name.Alias = alias[0]
	}
	return d
}

// Where 设置查询条件
//
// 参数：
//   - column: 列名或完整的条件语句
//   - w: 条件值或操作符和值
//
// 返回值：
//   - ModelTransfer: 模型转换接口
//
// 说明：
//   - Where("id", 1) => id = 1
//   - Where("id", ">", 1) => id > 1
//   - Where("id = 1 AND name = 'test'") => id = 1 AND name = 'test'
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
			frame.Errors(_tag, "[yorm][sql-where-func] w如果传递两个参数,第一个参数表示操作符,也就是说,只支持字符", nil)
			return d
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

// Join 设置连接
//
// 参数：
//   - table: 连接表
//   - condition: 连接条件
//   - link: 连接方式
//
// 返回值：
//   - ModelTransfer: 模型转换接口
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

// Field 设置查询字段
//
// 参数：
//   - column: 查询字段
//
// 返回值：
//   - ModelTransfer: 模型转换接口
func (d *StatementData) Field(column string) ModelTransfer {
	d.s.Fields = column
	return d
}

// Page 设置分页
//
// 参数：
//   - page: 页码
//   - size: 每页大小
//
// 返回值：
//   - ModelTransfer: 模型转换接口
func (d *StatementData) Page(page int, size int) ModelTransfer {
	if page < 1 {
		frame.Errors(_tag, "查询页码不能小于1", nil)
		return d
	}
	if size < 1 {
		frame.Errors(_tag, "查询条数不能小于1", nil)
		return d
	}
	d.s.Pages.Num = (page - 1) * size
	d.s.Pages.Size = size
	return d
}

// Limit 设置限制
//
// 参数：
//   - size: 限制大小
//
// 返回值：
//   - ModelTransfer: 模型转换接口
func (d *StatementData) Limit(size int) ModelTransfer {
	if size == 0 {
		frame.Errors(_tag, "查询条数不能小于1", nil)
		return d
	}
	d.s.Pages.Num = 0
	d.s.Pages.Size = size
	return d
}

// Order 设置排序
//
// 参数：
//   - column: 排序列
//   - sort: 排序方式
//
// 返回值：
//   - ModelTransfer: 模型转换接口
func (d *StatementData) Order(column string, sort string) ModelTransfer {
	o := Orders{
		column,
		sort,
	}
	d.s.Orders = append(d.s.Orders, o)
	return d
}

// Group 设置分组
//
// 参数：
//   - column: 分组列
//
// 返回值：
//   - ModelTransfer: 模型转换接口
func (d *StatementData) Group(column string) ModelTransfer {
	d.s.Groups = append(d.s.Groups, column)
	return d
}

// FetchSQL 设置返回SQL
//
// 参数：
//   - flag: 是否执行
//
// 返回值：
//   - ModelTransfer: 模型转换接口
func (d *StatementData) FetchSQL(flag ...bool) ModelTransfer {
	d.s.SQL = true
	if len(flag) > 0 && flag[0] {
		d.s.Exec = true
	}
	return d
}

// Find 查找单个
//
// 参数：
//   - fs: 格式化结构
//
// 返回值：
//   - QueryResult: 查询结果
func (d *StatementData) Find(fs ...FormatStruct) QueryResult {
	return d.s.d.driver.GetModel(d).Find(fs...)
}

// Select 查找多个
//
// 参数：
//   - fs: 格式化结构
//
// 返回值：
//   - QueryResults: 查询结果
func (d *StatementData) Select(fs ...FormatStruct) QueryResults {
	return d.s.d.driver.GetModel(d).Select(fs...)
}

// Update 更新
//
// 参数：
//   - i: 更新数据
//
// 返回值：
//   - ExecResult: 执行结果
func (d *StatementData) Update(i interface{}) ExecResult {
	return d.s.d.driver.GetModel(d).Update(i)
}

// Insert 插入
//
// 参数：
//   - i: 插入数据
//
// 返回值：
//   - ExecResult: 执行结果
func (d *StatementData) Insert(i interface{}) ExecResult {
	return d.s.d.driver.GetModel(d).Insert(i)
}

// Delete 删除
//
// 返回值：
//   - ExecResult: 执行结果
func (d *StatementData) Delete() ExecResult {
	return d.s.d.driver.GetModel(d).Delete()
}

// Db 获取数据库连接
//
// 返回值：
//   - *Db: 数据库连接
func (d *StatementData) Db() *Db {
	return d.s.d
}

// GetStatement 获取语句
//
// 返回值：
//   - *Statement: 语句
func (d *StatementData) GetStatement() *Statement {
	return &d.s
}

func (d *StatementData) Count() CountResult {
	return d.s.d.driver.GetModel(d).Count()
}
func (d *StatementData) Value(column string) ValueResult {
	return d.s.d.driver.GetModel(d).Value(column)
}
