package yorm

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
	"yiarce/core/frame"
	"yiarce/core/timing"

	_ "github.com/go-sql-driver/mysql"
)

var (
	driversMu sync.RWMutex
	d         = make(map[string]Driver)
)

const _tag = `yorm`

type Db struct {
	driver Driver
	conn   *sql.DB
	tx     *sql.Tx
}

type Driver interface {
	GetModel(ctx *StatementData) ImplTransfer
}

func Register(name string, driver Driver) {
	driversMu.Lock()
	defer driversMu.Unlock()
	if driver == nil {
		frame.Errors(_tag, "驱动不存在", nil)
		return
	}
	if _, dup := d[name]; dup {
		return
	}
	d[name] = driver
}

func Connect(config Config) (*Db, error) {
	switch config.Type {
	case `mysql`, `sqlite`:
		return ConnMysql(config)
	default:
		return nil, errors.New("暂不支持的数据库")
	}
}

// ConnMysql 创建MySQL数据库连接
//
// 参数：
//   - c: 数据库配置
//
// 返回值：
//   - *Db: 数据库连接对象
//   - error: 错误信息
func ConnMysql(c Config) (*Db, error) {
	var DbLink strings.Builder
	DbLink.WriteString(c.Username)
	DbLink.WriteString(":")
	DbLink.WriteString(c.Password)
	DbLink.WriteString("@tcp(")
	DbLink.WriteString(c.Host)
	DbLink.WriteString(":")
	DbLink.WriteString(c.Port)
	DbLink.WriteString(")/")
	DbLink.WriteString(c.Database)
	if c.Character != `` {
		DbLink.WriteString(`?charset=` + c.Character)
	}
	conn, err := sql.Open("mysql", DbLink.String())
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// 设置连接池参数
	conn.SetMaxOpenConns(150)
	conn.SetMaxIdleConns(20)
	conn.SetConnMaxLifetime(time.Hour) // 设置连接最大生命周期，避免连接泄漏

	// 测试连接
	if err := conn.Ping(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// 检查驱动是否初始化
	if d["mysql"] == nil {
		conn.Close()
		return nil, fmt.Errorf("mysql database driver not initialized, please import the corresponding driver package")
	}

	// 启动连接保活
	keepAlive(conn)

	return &Db{
		driver: d["mysql"],
		conn:   conn,
		tx:     nil,
	}, nil
}

// keepAlive 保持数据库连接活跃
//
// 参数：
//   - conn: 数据库连接
func keepAlive(conn *sql.DB) {
	timing.Anonymous(func() bool {
		err := conn.Ping()
		if err != nil {
			frame.Println(err.Error())
			return false
		}
		return true
	}, time.Second*60).Start()
}

// Query 执行查询并返回结果集
//
// 参数：
//   - query: SQL查询语句
//
// 返回值：
//   - *sql.Rows: 结果集
//   - error: 错误信息
func (d *Db) Query(query string) (*sql.Rows, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("[ERROR] YORM查询异常: %v\n", err)
		}
	}()
	return querySql(d.conn, d.tx, query)
}

// QueryMap 执行查询并返回映射格式的结果
//
// 参数：
//   - query: SQL查询语句
//
// 返回值：
//   - []map[string]string: 映射格式的结果集
//   - error: 错误信息
func (d *Db) QueryMap(query string) ([]map[string]string, error) {
	rows, err := d.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	column, _ := rows.Columns()
	values := make([]sql.RawBytes, len(column))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	// 预分配结果切片，减少动态扩容
	var res []map[string]string

	for rows.Next() {
		if err := rows.Scan(scanArgs...); err != nil {
			return nil, err
		}

		// 预分配映射，避免动态扩容
		rowMap := make(map[string]string, len(column))
		for i, col := range values {
			if col != nil {
				rowMap[column[i]] = string(col)
			}
		}
		res = append(res, rowMap)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

// Execute 执行SQL语句并返回结果
//
// 参数：
//   - query: SQL执行语句
//
// 返回值：
//   - int64: 最后插入的ID
//   - int64: 影响的行数
//   - error: 错误信息
func (d *Db) Execute(query string) (int64, int64, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("[ERROR] YORM执行异常: %v\n", err)
		}
	}()
	return executeSql(d.conn, d.tx, query)
}

// Table 设置查询表名
//
// 参数：
//   - name: 表名
//   - alias: 表别名
//
// 返回值：
//   - ModelTransfer: 模型转换接口
func (d *Db) Table(name string, alias ...string) ModelTransfer {
	s := sample
	s.Name.Name = name
	s.d = d
	l := len(alias)
	if l > 0 {
		s.Name.Alias = alias[0]
	}
	return &StatementData{s}
}

// Where 添加查询条件
//
// 参数：
//   - column: 列名或完整的WHERE条件
//   - w: 条件值或操作符和条件值
//
// 返回值：
//   - ModelTransfer: 模型转换接口
func (d *Db) Where(column string, w ...interface{}) ModelTransfer {
	s := sample
	s.d = d
	wh := Wheres{}
	l := len(w)
	switch l {
	case 1:
		wh.Column = column
		wh.Exp = "="
		wh.Content = checkType(w[0])
		s.Wheres = append(s.Wheres, wh)
	case 2:
		wh.Column = column
		exp, flag := w[0].(string)
		if !flag {
			frame.Errors(_tag, "if two arguments are passed, the first one represents the operator, that is, only characters are supported", nil)
			return &StatementData{s}
		}
		wh.Exp = exp
		wh.Content = checkType(w[1])
		s.Wheres = append(s.Wheres, wh)
	default:
		s.Wheres = append(s.Wheres, column)
	}
	return &StatementData{s}
}

// Join 添加连接查询
//
// 参数：
//   - table: 连接表名
//   - condition: 连接条件
//   - link: 连接方式（如INNER、LEFT、RIGHT等）
//
// 返回值：
//   - ModelTransfer: 模型转换接口
func (d *Db) Join(table string, condition string, link ...string) ModelTransfer {
	s := sample
	s.d = d
	j := Joins{}
	j.JoinTable = table
	j.JoinWhere = condition
	if len(link) > 0 {
		j.JoinWay = link[0]
	}
	return &StatementData{s}
}

// Page 设置分页参数
//
// 参数：
//   - page: 页码（从1开始）
//   - size: 每页大小
//
// 返回值：
//   - ModelTransfer: 模型转换接口
func (d *Db) Page(page int, size int) ModelTransfer {
	if page < 1 {
		frame.Errors(_tag, "page can't be less than 1", nil)
		s := sample
		s.d = d
		return &StatementData{s}
	}
	if size < 1 {
		frame.Errors(_tag, "size can't be less than 1", nil)
		s := sample
		s.d = d
		return &StatementData{s}
	}
	s := sample
	s.d = d
	s.Pages = Pages{
		Num:  (page - 1) * size,
		Size: size,
	}
	return &StatementData{s}
}

// Limit 设置查询限制
//
// 参数：
//   - size: 限制数量
//
// 返回值：
//   - ModelTransfer: 模型转换接口
func (d *Db) Limit(size int) ModelTransfer {
	if size == 0 {
		frame.Errors(_tag, "the size cannot be less than 1", nil)
		s := sample
		s.d = d
		return &StatementData{s}
	}
	s := sample
	s.d = d
	s.Pages = Pages{
		Num:  0,
		Size: size,
	}
	return &StatementData{s}
}

// Filed 设置查询字段
//
// 参数：
//   - columns: 查询字段
//
// 返回值：
//   - ModelTransfer: 模型转换接口
func (d *Db) Filed(columns string) ModelTransfer {
	s := sample
	s.d = d
	s.Fields = columns
	return &StatementData{s}
}

// Group 设置分组字段
//
// 参数：
//   - column: 分组字段
//
// 返回值：
//   - ModelTransfer: 模型转换接口
func (d *Db) Group(column string) ModelTransfer {
	s := sample
	s.d = d
	var g []string
	g = append(g, column)
	s.Groups = g
	return &StatementData{s}
}

// Order 设置排序字段
//
// 参数：
//   - column: 排序字段
//   - sort: 排序方式（如ASC、DESC）
//
// 返回值：
//   - ModelTransfer: 模型转换接口
func (d *Db) Order(column string, sort string) ModelTransfer {
	s := sample
	s.d = d
	s.Orders = []Orders{{
		Column: column,
		Sort:   sort,
	}}
	return &StatementData{s}
}

// Close 关闭数据库连接
//
// 返回值：
//   - error: 错误信息
func (d *Db) Close() error {
	if d.conn != nil {
		return d.conn.Close()
	}
	return nil
}

// Begin 开始事务
//
// 返回值：
//   - error: 错误信息
func (d *Db) Begin() (*Db, error) {
	tx, err := d.conn.Begin()
	if err != nil {
		return nil, err
	}
	return &Db{
		conn: d.conn,
		tx:   tx,
	}, nil
}

// Commit 提交事务
//
// 返回值：
//   - error: 错误信息
func (d *Db) Commit() error {
	if d.tx == nil {
		return errors.New("no transaction started")
	}
	err := d.tx.Commit()
	if err != nil {
		return err
	}
	d.tx = nil
	return nil
}

// Rollback 回滚事务
//
// 返回值：
//   - error: 错误信息
func (d *Db) Rollback() error {
	if d.tx == nil {
		return errors.New("no transaction started")
	}
	err := d.tx.Rollback()
	if err != nil {
		return err
	}
	d.tx = nil
	return nil
}

func querySql(db *sql.DB, tx *sql.Tx, s string) (*sql.Rows, error) {
	if tx != nil {
		return tx.Query(s)
	}
	return db.Query(s)
}

func executeSql(db *sql.DB, tx *sql.Tx, query string) (int64, int64, error) {
	var result sql.Result
	var err error
	if tx != nil {
		result, err = tx.Exec(query)
	} else {
		result, err = db.Exec(query)
	}
	if err != nil {
		return 0, 0, err
	}
	rowNum, err := result.RowsAffected()
	if err != nil {
		return 0, 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, rowNum, nil
	}
	return id, rowNum, nil
}
