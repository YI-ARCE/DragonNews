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
	FetchSql(flag ...bool) ModelTransfer
	Find(fs ...FormatStruct) QueryResult
	Select(fs ...FormatStruct) QueryResults
	Insert(i interface{}) ExecResult
	Update(i interface{}) ExecResult
	Delete() ExecResult
}

type ImplTransfer interface {
	Find(fs ...FormatStruct) QueryResult
	Select(fs ...FormatStruct) QueryResults
	Insert(i interface{}) ExecResult
	Update(i interface{}) ExecResult
	Delete() ExecResult
}

type Config struct {
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
	Format(p FormatStruct)
	Sql() string
	Concat(key string, flag bool) string
}

type QueryResult interface {
	Result() map[string]string
	Err() error
	Format(p FormatStruct)
	Sql() string
}

type ExecResult interface {
	Id() int64
	Num() int64
	Err() error
	Sql() string
}
