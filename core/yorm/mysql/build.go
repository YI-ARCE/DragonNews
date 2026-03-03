package mysql

import (
	"reflect"
	"strconv"
	"strings"
	"yiarce/core/frame"
	"yiarce/core/yorm"
)

const _tag = `mysql`

func checkType(i interface{}, f ...bool) string {
	r := reflect.ValueOf(i)
	switch r.Type().Kind() {
	case reflect.String:
		i := r.String()
		if len(i) > 9 && i[:10] == `[raw]__dn:` {
			return i[10:]
		}
		if len(f) > 0 {
			return `"` + strings.ReplaceAll(strings.ReplaceAll(i, `\`, `\\`), `"`, `\"`) + `"`
		}
		return strings.ReplaceAll(strings.ReplaceAll(i, `\`, `\\`), `"`, `\"`)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(r.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(r.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(r.Float(), 'f', -1, 64)
	default:
		frame.Errors(_tag, "不支持的解包类型: "+r.Kind().String(), nil)
		return ""
	}
}

func parseAlias(n yorm.Name) string {
	if n.Alias != "" {
		return strBuild(" ", "`", n.Alias, "`")
	}
	return ""
}

func parseJoins(js []yorm.Joins) string {
	sql := ""
	for _, j := range js {
		way := "INNER"
		if j.JoinWay != "" {
			way = j.JoinWay
		}
		sql = strJoin(sql, " ", way, " JOIN ")
		if strings.Contains(j.JoinTable, " ") && strings.Count(j.JoinTable, " ") == 1 {
			str := strings.Split(j.JoinTable, " ")
			sql = strJoin(sql, str[0], " `", str[1], "`")
		} else {
			sql = strJoin(sql, j.JoinTable)
		}
		sql = strJoin(sql, " ON ( ", j.JoinWhere, " ) ")
	}
	return sql
}

func parseWhere(ws []interface{}) string {
	sql := ""
	for _, v := range ws {
		str, flag := v.(string)
		if flag {
			sql = strJoin(sql, " ( ", str, " ) AND")
		} else {
			w := v.(yorm.Wheres)
			if strings.Contains(w.Column, ".") {
				cl := strings.Split(w.Column, ".")
				sql = strJoin(sql, " ( ", cl[0], ".`", cl[1], "` ")
			} else {
				sql = strJoin(sql, " ( ", w.Column, " ")
			}
			if w.Exp != "" {
				sql = strBuild(sql, w.Exp, " ", w.Content)
			} else {
				sql = strBuild(sql, " ", w.Content)
			}
			sql = strBuild(sql, " ) AND")
		}
	}

	if sql != "" {
		sql = strBuild(" WHERE", sql[:len(sql)-4])
	}
	return sql
}

func parseGroup(gs []string) string {
	sql := ""
	for _, g := range gs {
		if len(g) > 15 {
			sql = strJoin(sql, " GROUP BY ", g)
		} else {
			if strings.Contains(g, ".") {
				str := strings.Split(g, ".")
				sql = strBuild(sql, str[0], "`", str[1], "`")
			} else {
				sql = strJoin(sql, "GROUP BY `", g, "`")
			}
		}
	}
	return sql
}

func parseOrder(os []yorm.Orders) string {
	if len(os) > 0 {
		sql := " ORDER BY"
		for _, o := range os {
			if strings.Contains(o.Column, ".") {
				str := strings.Split(o.Column, ".")
				sql = strBuild(sql, " `", str[0], "`.`", str[1], "`")
			} else {
				sql = strBuild(sql, " `", o.Column, "` ")
			}
			if o.Sort != "" {
				sql = strBuild(sql, ` `, o.Sort, ` `)
			}
			sql = strBuild(sql, `,`)
		}
		return sql[:len(sql)-1]
	}
	return ``
}

func parsePage(p yorm.Pages) string {
	if p.Size == 0 {
		return ""
	}
	if p.Num > 0 {
		return strBuild(" LIMIT ", strconv.Itoa(p.Num*p.Size), ",", strconv.Itoa(p.Size))
	} else {
		return strBuild(" LIMIT ", strconv.Itoa(p.Size))
	}
}

func query(c *yorm.Statement) string {
	sql := "SELECT "
	sql = strBuild(sql, c.Fields, " FROM ", c.Name.Name)
	sql = strBuild(sql, parseAlias(c.Name))
	sql = strJoin(sql, parseJoins(c.Joins))
	sql = strJoin(sql, parseWhere(c.Wheres))
	sql = strBuild(sql, parseGroup(c.Groups))
	sql = strBuild(sql, parseOrder(c.Orders))
	sql = strBuild(sql, parsePage(c.Pages))
	return sql
}

func update(c *yorm.Statement, i interface{}) string {
	sql := strBuild("UPDATE ", c.Name.Name)
	sql = strBuild(sql, parseAlias(c.Name))
	sql = strJoin(sql, parseJoins(c.Joins), " SET ")
	t := reflect.TypeOf(i)
	r := reflect.ValueOf(i)
	switch t.Kind() {
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			value := r.Field(i)
			// 跳过未导出的字段
			if field.PkgPath != "" {
				continue
			}
			// 获取字段名，优先使用json标签
			fieldName := field.Name
			if tag := field.Tag.Get("json"); tag != "" {
				fieldName = strings.Split(tag, ",")[0]
			}
			// 跳过空值字段
			if value.IsZero() {
				continue
			}
			sql = strBuild(sql, fieldName, " = ", checkType(value.Interface(), true), ",")
		}
	case reflect.Map:
		m := r.MapRange()
		for m.Next() {
			if m.Key().Kind() != reflect.String {
				frame.Errors(_tag, `map键类型必须为string`, nil)
				return ""
			}
			sql = strBuild(sql, m.Key().String(), " = ", checkType(m.Value().Interface(), true), ",")
		}
	case reflect.String:
		i := i.(string)
		i = strings.ReplaceAll(strings.ReplaceAll(i, `\`, `\\`), `"`, `\"`)
		if len(i) > 9 && i[:10] == `[raw]__dn:` {
			sql = strBuild(sql, i[10:], ",")
		} else {
			sql = strBuild(sql, i, ",")
		}
	default:
		frame.Errors(_tag, `仅支持map与struct类型`, nil)
		return ""
	}
	// 移除最后的逗号
	if strings.HasSuffix(sql, ",") {
		sql = sql[:len(sql)-1]
	}
	sql = strJoin(sql, parseWhere(c.Wheres))
	return sql
}

func exec(c *yorm.Statement, i interface{}) string {
	sql := strBuild("INSERT INTO ", c.Name.Name)
	t := reflect.TypeOf(i)
	r := reflect.ValueOf(i)
	switch t.Kind() {
	case reflect.Struct:
		keyStr := ""
		valStr := ""
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			value := r.Field(i)
			// 跳过未导出的字段
			if field.PkgPath != "" {
				continue
			}
			// 跳过空值字段
			if value.IsZero() {
				continue
			}
			// 获取字段名，优先使用json标签
			fieldName := field.Name
			if tag := field.Tag.Get("json"); tag != "" {
				fieldName = strings.Split(tag, ",")[0]
			}
			keyStr = strBuild(keyStr, "`", fieldName, "`", ",")
			valStr = strBuild(valStr, checkType(value.Interface(), true), ",")
		}
		// 移除最后的逗号
		if strings.HasSuffix(keyStr, ",") {
			keyStr = keyStr[:len(keyStr)-1]
		}
		if strings.HasSuffix(valStr, ",") {
			valStr = valStr[:len(valStr)-1]
		}
		sql = strJoin(sql, " (", keyStr, ") VALUES (", valStr, ")")
	case reflect.Map:
		m := r.MapRange()

		keyStr := ""
		valStr := ""
		for m.Next() {
			if m.Key().Kind() != reflect.String {
				frame.Errors(_tag, `map键类型必须为string`, nil)
				return ""
			}
			keyStr = strBuild(keyStr, "`", m.Key().String(), "`,")
			valStr = strBuild(valStr, checkType(m.Value().Interface(), true), ",")
		}
		keyStr = keyStr[0 : len(keyStr)-1]
		valStr = valStr[0 : len(valStr)-1]
		sql = strJoin(sql, " (", keyStr, ") VALUES (", valStr, ")")
	default:
		frame.Errors(_tag, `仅支持map与struct类型`, nil)
		return ""
	}
	return sql
}

func remove(c *yorm.Statement) string {
	sql := "DELETE "
	if c.Fields != "*" {
		sql = strBuild(sql, c.Fields, " FROM ", c.Name.Name)
	} else {
		sql = strBuild(sql, "FROM ", c.Name.Name)
	}
	sql = strBuild(sql, parseAlias(c.Name))
	sql = strJoin(sql, parseJoins(c.Joins))
	sql = strJoin(sql, parseWhere(c.Wheres))
	sql = strJoin(sql, parsePage(c.Pages))
	return sql
}

func strJoin(str string, s ...string) string {
	bs := []string{str}
	for _, v := range s {
		bs = append(bs, v)
	}
	return strings.Join(bs, "")
}

func strBuild(str string, s ...string) string {
	b := strings.Builder{}
	b.WriteString(str)
	for _, v := range s {
		b.WriteString(v)
	}
	return b.String()
}
