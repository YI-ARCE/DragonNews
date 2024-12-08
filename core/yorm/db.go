package yorm

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	"sync"
	"time"
	"yiarce/core/frame"
	"yiarce/core/timing"
)

var (
	driversMu sync.RWMutex
	d         = make(map[string]Driver)
)

type Db struct {
	driver Driver
	conn   *sql.DB
}

type Driver interface {
	GetModel(ctx *StatementData) ImplTransfer
}

func Register(name string, driver Driver) {
	driversMu.Lock()
	defer driversMu.Unlock()
	if driver == nil {
		panic("Sql: Register driver is nil")
	}
	if _, dup := d[name]; dup {
		panic("Sql: Register called twice for driver " + name)
	}
	d[name] = driver
}

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
	conn.SetMaxOpenConns(150)
	conn.SetMaxIdleConns(20)
	if err != nil {
		return nil, err
	} else {
		err := conn.Ping()
		if err != nil {
			return nil, err
		}
		if d["mysql"] == nil {
			panic(`mysql数据库驱动未初始化,请尝试引入对应的驱动包`)
		}
		keepAlive(conn)
		return &Db{
			d["mysql"],
			conn,
		}, nil
	}
}

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

func (d *Db) Query(query string) (*sql.Rows, error) {
	return querySql(d.conn, query)
}

func (d *Db) QueryMap(query string) ([]map[string]string, error) {
	rows, err := d.Query(query)
	if err != nil {
		return nil, err
	}
	var res []map[string]string
	column, _ := rows.Columns()
	values := make([]sql.RawBytes, len(column))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	for rows.Next() {
		_ = rows.Scan(scanArgs...)
		rowMap := make(map[string]string)
		var value string
		for i, col := range values {
			if col != nil {
				value = string(col)
				rowMap[column[i]] = value
			}
		}
		res = append(res, rowMap)
	}
	rows.Close()
	return res, nil
}

func (d *Db) Execute(query string) (int64, int64, error) {
	return executeSql(d.conn, query)
}

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
			panic("if two arguments are passed, the first one represents the operator, that is, only characters are supported")
		}
		wh.Exp = exp
		wh.Content = checkType(w[0])
		s.Wheres = append(s.Wheres, wh)
	default:
		s.Wheres = append(s.Wheres, column)
	}
	return &StatementData{s}
}

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

func (d *Db) Page(page int, size int) ModelTransfer {
	if page < 1 {
		panic("page can't be less than 1")
	}
	s := sample
	s.Pages = Pages{
		Num:  page - 1,
		Size: size * page,
	}
	return &StatementData{s}
}

func (d *Db) Limit(size int) ModelTransfer {
	if size == 0 {
		panic("the size cannot be less than 1")
	}
	s := sample
	s.d = d
	s.Pages = Pages{
		Num:  0,
		Size: size,
	}
	return &StatementData{s}
}

func (d *Db) Filed(columns string) ModelTransfer {
	s := sample
	s.d = d
	s.Fields = columns
	return &StatementData{s}
}

func (d *Db) Group(column string) ModelTransfer {
	s := sample
	s.d = d
	var g []string
	g = append(g, column)
	s.Groups = g
	return &StatementData{s}
}

func (d *Db) Order(column string, sort string) ModelTransfer {
	s := sample
	s.d = d
	s.Orders = []Orders{{
		Column: column,
		Sort:   sort,
	}}
	return &StatementData{s}
}

func (d *Db) Close() {
	d.conn.Close()
}

func querySql(db *sql.DB, s string) (*sql.Rows, error) {
	return db.Query(s)
}

func executeSql(bx tr, query string) (int64, int64, error) {
	result, err := bx.Exec(query)
	if err != nil {
		return 0, 0, err
	}
	rowNum, err := result.RowsAffected()
	if err != nil {
		return 0, 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, 0, err
	}
	return id, rowNum, nil
}
