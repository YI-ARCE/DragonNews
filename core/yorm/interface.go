package yorm

import "database/sql"

type tr interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type ModelTransfer interface {
	Table(name string, alias ...string) ModelTransfer
	Where(column string, w ...interface{}) ModelTransfer
	Join(table string, condition string, link ...string) ModelTransfer
	Field(column string) ModelTransfer
	Page(page int, size int) ModelTransfer
	Limit(size int) ModelTransfer
	Order(column string, sort string) ModelTransfer
	Group(column string) ModelTransfer
	FetchSQL(flag ...bool) ModelTransfer
	Find(fs ...FormatStruct) QueryResult
	Select(fs ...FormatStruct) QueryResults
	Insert(i interface{}) ExecResult
	Update(i interface{}) ExecResult
	Count() CountResult
	Value(column string) ValueResult
	Sum(column string) SumResult
	Delete() ExecResult
}

type ImplTransfer interface {
	Find(fs ...FormatStruct) QueryResult
	Select(fs ...FormatStruct) QueryResults
	Insert(i interface{}) ExecResult
	Update(i interface{}) ExecResult
	Count() CountResult
	Value(column string) ValueResult
	Sum(column string) SumResult
	Delete() ExecResult
}

type Config struct {
	Type      string `yaml:"type"`
	Host      string `yaml:"host"`
	Port      string `yaml:"port"`
	Database  string `yaml:"database"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	Character string `yaml:"character"`
}

type QueryResults interface {
	Result() []map[string]string
	Err() error
	Sql() string
	Concat(key string, flag bool) string
}

type QueryResult interface {
	Result() map[string]string
	Err() error
	Sql() string
}

type ExecResult interface {
	Id() int64
	Num() int64
	Err() error
	Sql() string
}

type CountResult interface {
	Num() int64
	Err() error
	Sql() string
}

type SumResult interface {
	Num() int64
	Err() error
	Sql() string
}

type ValueResult interface {
	String() string
	Err() error
	Sql() string
	// Float 获取浮点值
	Float64() float64
	Float() float32
	Int() int
	Int64() int64
	IsNil() bool
	// Date 只有值为时间戳且类型为int时才有效,获取日期值
	Date(format ...string) string
}
