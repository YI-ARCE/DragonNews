package mysql

import (
	"database/sql"
	"yiarce/core/yorm"
)

type er struct {
	id  int64
	num int64
	err error
	sql string
}

type qr struct {
	err    error
	data   interface{}
	sql    string
	result *sql.Rows
}

type qfr struct {
	result *sql.Rows
	err    error
	data   interface{}
	sql    string
}

func (e er) Id() int64 {
	return e.id
}

func (e er) Num() int64 {
	return e.num
}

func (e er) Err() error {
	return e.err
}

func (e er) Sql() string {
	return e.sql
}

func (q *qr) Result() []map[string]string {
	return toMap(q.result)
}

func (q *qr) Err() error {
	return q.err
}

func (q *qr) Format(p yorm.FormatStruct) {

}

func (q *qr) Sql() string {
	return q.sql
}

func (q *qr) Concat(key string, flag bool) string {
	str := ``
	t1 := ``
	t2 := ``
	if flag {
		t1, t2 = `'`, `'`
	}
	for _, v := range q.Result() {
		str += t1 + v[key] + t2 + `,`
	}
	l := len(str)
	if l > 0 {
		return str[:l-1]
	}
	return str
}

func (q *qfr) Result() map[string]string {
	if q.result != nil {
		maps := toMap(q.result)
		if len(maps) > 0 {
			return maps[0]
		}
	}
	return map[string]string{}
}

func (q *qfr) Err() error {
	return q.err
}

func (q *qfr) Format(p yorm.FormatStruct) {
	//t := reflect.TypeOf(p)
	//v := reflect.ValueOf(p)
	//l := t.NumField()
	//for i := 0; i < l; i++ {
	//
	//}
}

func (q *qfr) Sql() string {
	return q.sql
}

func toMap(result *sql.Rows) []map[string]string {
	res := make([]map[string]string, 0)
	column, _ := result.Columns()
	values := make([]sql.RawBytes, len(column))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	for result.Next() {
		_ = result.Scan(scanArgs...)
		rowMap := make(map[string]string)
		var value string
		for i, col := range values {
			if col != nil {
				value = string(col)
				rowMap[column[i]] = value
			} else {
				rowMap[column[i]] = ``
			}
		}
		res = append(res, rowMap)
	}
	return res
}
