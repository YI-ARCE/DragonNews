package mysql

import (
	"database/sql"
	"strconv"
	"yiarce/core/date"
)

func toMap(result *sql.Rows) []map[string]string {
	res := make([]map[string]string, 0)
	if result == nil {
		return res
	}
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

func toFloat64(v string) float64 {
	f, _ := strconv.ParseFloat(v, 64)
	return f
}

func toInt64(v string) int64 {
	i, _ := strconv.ParseInt(v, 10, 64)
	return i
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

func (q qr) Result() []map[string]string {
	return toMap(q.result)
}

func (q qr) Err() error {
	return q.err
}

func (q qr) Sql() string {
	return q.sql
}

func (q qr) Concat(key string, flag bool) string {
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

func (q qfr) Result() map[string]string {
	if q.result != nil {
		maps := toMap(q.result)
		if len(maps) > 0 {
			return maps[0]
		}
	}
	return map[string]string{}
}

func (q qfr) Err() error {
	return q.err
}

func (q qfr) Sql() string {
	return q.sql
}

type cfr struct {
	err error
	num int64
	sql string
}

func (c cfr) Num() int64 {
	return c.num
}

func (c cfr) Err() error {
	return c.err
}

func (c cfr) Sql() string {
	return c.sql
}

type sfr struct {
	err error
	num int64
	sql string
}

func (s sfr) Num() int64 {
	return s.num
}

func (s sfr) Err() error {
	return s.err
}

func (s sfr) Sql() string {
	return s.sql
}

type vfr struct {
	err   error
	value string
	sql   string
	nil   bool
}

func (v vfr) Err() error {
	return v.err
}

func (v vfr) Sql() string {
	return v.sql
}

func (v vfr) Float64() float64 {
	return toFloat64(v.value)
}

func (v vfr) Float() float32 {
	return float32(toFloat64(v.value))
}

func (v vfr) Int() int {
	return int(toInt64(v.value))
}

func (v vfr) String() string {
	return v.value
}

func (v vfr) Bool() bool {
	return string(v.value) == `true`
}

func (v vfr) Int64() int64 {
	return toInt64(v.value)
}

func (v vfr) IsNil() bool {
	return v.nil
}

func (v vfr) Date(format ...string) string {
	if len(format) > 0 {
		return date.Time(toInt64(v.value)).Custom(format[0])
	}
	return date.Time(toInt64(v.value)).Custom()
}
