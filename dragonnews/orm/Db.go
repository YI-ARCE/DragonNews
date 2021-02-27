package Db

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strconv"
	"strings"
	"time"
)

type Sql struct {
	Type     string
	Host     string
	Port     string
	Database string
	Username string
	Password string
}

var Conn *sql.DB

type query struct {
	table  string
	where  string
	field  string
	order  string
	group  string
	join   string
	limit  string
	page   string
	sql    bool
	types  string
	data   map[string]string
	lastId int
	bx     *sql.Tx
}

func Init(SqlConfig Sql) {
	var DbLink strings.Builder
	DbLink.WriteString(SqlConfig.Username)
	DbLink.WriteString(":")
	DbLink.WriteString(SqlConfig.Password)
	DbLink.WriteString("@tcp(")
	DbLink.WriteString(SqlConfig.Host)
	DbLink.WriteString(":")
	DbLink.WriteString(SqlConfig.Port)
	DbLink.WriteString(")/")
	DbLink.WriteString(SqlConfig.Database)
	var err error
	Conn, err = sql.Open(SqlConfig.Type, DbLink.String())
	if err != nil {
		log.Print(err)
	}
	err = Conn.Ping()
	if err != nil {
		fmt.Println(err)
	}
	Conn.SetMaxOpenConns(10)
	Conn.SetMaxIdleConns(1)
	Conn.SetConnMaxLifetime(time.Second * 3600)
}

//自定义查询,若查询的结果只返回一个则索引[0]即可
func Query(sql string) map[int]map[string]string {
	return checkSql(sql)
}

//自定义写入,结果1返回影响的条数,若为插入时结果2会返回本次插入的ID,更新操作不需要识别ID
func Exec(sql string) (int, int) {
	_, num, id := executeSql(sql)
	return num, id
}

//开启事务
func BeginTranse() (*query, error) {
	tx, err := Conn.Begin()
	if err != nil {
		return nil, err
	}
	query := Table("")
	query.bx = tx
	return query, nil
}

//提交事务
func (query *query) Commit() error {
	if query.bx == nil {
		return errors.New("不存在的事务!")
	}
	err := query.bx.Commit()
	if err != nil {
		return err
	}
	return nil
}

//回滚
func (query *query) RollBack() error {
	if query.bx == nil {
		return errors.New("不存在的事务!")
	}
	err := query.bx.Rollback()
	if err != nil {
		return err
	}
	return nil
}

//表
func Table(table string) *query {
	query := query{}
	query.table = table
	query.field = "*"
	query.sql = false
	return &query
}

//表
func (query *query) Table(table string) *query {
	query.table = table
	query.field = "*"
	query.sql = false
	return query
}

//条件
func (query *query) Where(where string) *query {
	var builder strings.Builder
	if query.where != "" {
		builder.WriteString(query.where)
		builder.WriteString(" AND (")
	} else {
		builder.WriteString(" (")
	}
	builder.WriteString(where)
	builder.WriteString(")")

	query.where = builder.String()
	return query
}

//需要字段
func (query *query) Field(field string) *query {
	var str []string
	var builder strings.Builder
	if query.field != "*" {
		builder.WriteString(query.field)
		builder.WriteString(",")
	}
	flag := strings.Contains(field, ",")
	if flag {
		str = strings.Split(field, ",")
		count := len(str)
		for key, val := range str {
			if strings.Contains(val, "(") || strings.Contains(val, ")") {
				builder.WriteString(val)
			} else if strings.Contains(val, ".") {
				column := strings.Split(val, ".")
				builder.WriteString("`")
				builder.WriteString(column[0])
				builder.WriteString("`.`")
				builder.WriteString(column[1])
				builder.WriteString("`")
			} else {
				builder.WriteString("`")
				builder.WriteString(val)
				builder.WriteString("`")
			}
			if key+1 < count {
				builder.WriteString(",")
			}
		}
	} else {
		if strings.Contains(field, ".") {
			column := strings.Split(field, ".")
			builder.WriteString("`")
			builder.WriteString(column[0])
			builder.WriteString("`.`")
			builder.WriteString(column[1])
			builder.WriteString("`")
		} else {
			builder.WriteString("`")
			builder.WriteString(field)
			builder.WriteString("`")
		}
	}
	query.field = builder.String()
	return query
}

//连表查询
func (query *query) Join(table string, option ...string) *query {
	var builder strings.Builder
	if query.join != "" {
		builder.WriteString(query.join)
		builder.WriteString(" ")
	}
	if len(option) < 2 {
		builder.WriteString("INNER JOIN ")
	} else {
		builder.WriteString(option[1])
		builder.WriteString(" JOIN ")
	}
	if strings.Contains(table, " ") && strings.Count(table, " ") < 1 {
		str := strings.Split(table, " ")
		builder.WriteString("`")
		builder.WriteString(str[0])
		builder.WriteString("` `")
		builder.WriteString(str[1])
		builder.WriteString("`")
	} else {
		builder.WriteString(table)
	}
	builder.WriteString(" ON ")
	builder.WriteString(option[0])

	query.join = builder.String()
	return query
}

func (query *query) Group(group string) *query {
	query.group = group
	return query
}

//排序
func (query *query) Order(order string, types ...string) *query {
	if len(types) > 0 {
		query.order = order + " " + types[0]
	} else {
		query.order = order
	}
	return query
}

func (query *query) Limit(limit int) *query {
	query.limit = strconv.Itoa(limit)
	return query
}

func (query *query) Page(page int, number int) *query {
	var builder strings.Builder
	if page < 1 {
		panic("页码不能小于1")
	}
	builder.WriteString(strconv.Itoa(page - 1))
	builder.WriteString(",")
	builder.WriteString(strconv.Itoa(number))
	query.page = builder.String()
	return query
}

//查询单个结果
func (query *query) Find() (map[string]string, string) {
	query.types = "check"
	sqlQuery := query.createQuery("LIMIT 1")
	if query.sql {
		return map[string]string{}, sqlQuery
	}
	result := checkSql(sqlQuery, query.bx)
	return result[0], ""

}

//查询结果集
func (query *query) Select() (map[int]map[string]string, string) {
	query.types = "check"
	sqlQuery := query.createQuery("")
	if query.sql {
		return map[int]map[string]string{}, sqlQuery
	}
	result := checkSql(sqlQuery, query.bx)
	return result, ""
}

func (query *query) Update(data map[string]string) (int, string) {
	query.types = "update"
	query.data = data
	sqlQuery := query.createQuery("")
	if query.sql {
		return 0, sqlQuery
	}
	_, num, _ := executeSql(sqlQuery, query.bx)
	return num, ""
}

func (query *query) Insert(data map[string]string) (bool, string) {
	query.types = "insert"
	query.data = data
	sqlQuery := query.createQuery("")
	if query.sql {
		return false, sqlQuery
	}
	flag, _, _ := executeSql(sqlQuery, query.bx)
	return flag, ""
}

func (query *query) GetLastId() int {
	return query.lastId
}

func (query *query) Delete() (int, string) {
	query.types = "delete"
	sqlQuery := query.createQuery("")
	if query.sql {
		return 0, sqlQuery
	}
	_, num, _ := executeSql(sqlQuery, query.bx)
	return num, ""
}

//生成SQL语句
func (query *query) Fet() *query {
	query.sql = true
	return query
}

//生成SQL语句
func (query *query) createQuery(limits string) string {
	//"SELECT * FROM `crm_member` INNER JOIN `bbb` `b` ON `b`.`ba`=`c`.`ca` WHERE ( status != 5 ) GROUP BY `status` ORDER BY `bbb` DESC LIMIT 1"
	var builder strings.Builder
	switch query.types {
	case "check":
		builder.WriteString("SELECT ")
		builder.WriteString(query.field)
		if strings.Contains(query.table, " ") {
			builder.WriteString(" FROM ")
			builder.WriteString(query.table)
		} else {
			builder.WriteString(" FROM `")
			builder.WriteString(query.table)
			builder.WriteString("`")
		}
		break
	case "update":
		if strings.Contains(query.table, " ") {
			builder.WriteString("UPDATE ")
			builder.WriteString(query.table)
		} else {
			builder.WriteString("UPDATE `")
			builder.WriteString(query.table)
			builder.WriteString("`")
		}
		if query.join != "" {
			builder.WriteString(" ")
			builder.WriteString(query.join)
		}
		builder.WriteString(" SET ")
		count := len(query.data)
		i := 1
		for key, val := range query.data {
			builder.WriteString("`")
			builder.WriteString(key)
			builder.WriteString("`")
			builder.WriteString(" = ")
			_, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				builder.WriteString(`"`)
				builder.WriteString(val)
				builder.WriteString(`"`)
			} else {
				builder.WriteString(val)
			}
			if i < count {
				builder.WriteString(",")
				i++
			}
		}
		break
	case "insert":
		if strings.Contains(query.table, " ") {
			builder.WriteString("INSERT INTO ")
			builder.WriteString(query.table)
		} else {
			builder.WriteString("INSERT INTO `")
			builder.WriteString(query.table)
			builder.WriteString("`")
		}
		builder.WriteString(" (")
		count := len(query.data)
		i := 1
		for key := range query.data {
			builder.WriteString("`")
			builder.WriteString(key)
			builder.WriteString("`")
			if i < count {
				builder.WriteString(",")
				i++
			}
		}
		builder.WriteString(") VALUES (")
		i = 1
		for _, val := range query.data {
			builder.WriteString("\"")
			builder.WriteString(val)
			builder.WriteString("\"")
			if i < count {
				builder.WriteString(",")
				i++
			}
		}
		builder.WriteString(")")
		break
	case "delete":
		if strings.Contains(query.table, " ") {
			builder.WriteString("DELETE FROM ")
			builder.WriteString(query.table)
		} else {
			builder.WriteString("DELETE FROM `")
			builder.WriteString(query.table)
			builder.WriteString("`")
		}
		if query.join != "" {
			builder.WriteString(" ")
			builder.WriteString(query.join)
		}
		break
	default:
		panic("No ending method selected")
	}
	if query.join != "" && query.types == "check" {
		builder.WriteString(" ")
		builder.WriteString(query.join)
	}
	if query.where != "" && query.types != "insert" {
		builder.WriteString(" WHERE")
		builder.WriteString(query.where)
	}
	if query.group != "" && query.types != "insert" {
		builder.WriteString(" GROUP BY ")
		builder.WriteString(query.group)
	}
	if query.order != "" && query.types != "insert" {
		builder.WriteString(" ORDER BY ")
		builder.WriteString(query.order)
	}
	builder.WriteString(" ")
	if limits != "" {
		builder.WriteString(limits)
	} else {
		if query.page != "" {
			builder.WriteString("LIMIT ")
			builder.WriteString(query.page)
		} else if query.limit != "" {
			builder.WriteString("LIMIT ")
			builder.WriteString(query.limit)
		}
	}

	return builder.String()
}

//查询SQL并返回结果
func checkSql(query string, bx ...*sql.Tx) map[int]map[string]string {
	result := make(map[int]map[string]string)
	var Row *sql.Rows
	var err error
	if len(bx) > 0 && bx[0] != nil {
		Row, err = bx[0].Query(query)
	} else {
		Row, err = Conn.Query(query)
	}
	if err != nil {
		return map[int]map[string]string{}
	}
	column, _ := Row.Columns()
	values := make([]sql.RawBytes, len(column))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	num := 0
	for Row.Next() {
		_ = Row.Scan(scanArgs...)
		rowMap := make(map[string]string)
		var value string
		for i, col := range values {
			if col != nil {
				value = string(col)
				rowMap[column[i]] = value
			}
		}
		result[num] = rowMap
		num++
	}
	return result
}

func executeSql(query string, bx ...*sql.Tx) (bool, int, int) {
	var result sql.Result
	var err error
	if len(bx) > 0 && bx[0] != nil {
		result, err = bx[0].Exec(query)
	} else {
		result, err = Conn.Exec(query)
	}
	flag := true
	if err != nil {
		panic(err)
	}
	rowNum, err := result.RowsAffected()
	if err != nil {
		panic(err)
	}
	lastId, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}
	if rowNum <= 0 {
		flag = false
	}
	return flag, int(rowNum), int(lastId)
}
