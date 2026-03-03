package mysql

import (
	sql2 "database/sql"
	"strings"
	"yiarce/core/yorm"
)

func init() {
	yorm.Register("mysql", Transfer{})
}

type Count struct {
	Count int64 `json:"dn_count"`
}

type Sum struct {
	Sum int64 `json:"dn_sum"`
}

type Model struct {
	s *yorm.StatementData
}

type Transfer struct {
}

type Db interface {
	Db() *yorm.Db
}

func (t Transfer) GetModel(s *yorm.StatementData) yorm.ImplTransfer {
	return Model{s}
}

func (m Model) Find(fs ...yorm.FormatStruct) yorm.QueryResult {
	m.s.Page(1, 1)
	sql := query(m.s.GetStatement())
	q := qfr{}
	if m.s.GetStatement().SQL {
		q.sql = sql
	}
	if m.s.GetStatement().Exec {
		return q
	}
	r, e := m.s.Db().Query(sql)
	if e != nil {
		q.err = e
		return q
	}
	if len(fs) > 0 {
		parseResponseData(r, fs[0])
		r.Close()
	} else {
		q.result = r
	}
	return q
}

func (m Model) Select(fs ...yorm.FormatStruct) yorm.QueryResults {
	sql := query(m.s.GetStatement())
	q := qr{}
	if m.s.GetStatement().SQL {
		q.sql = sql
	}
	if m.s.GetStatement().Exec {
		return q
	}
	r, e := m.s.Db().Query(sql)
	if e != nil {
		q.err = e
		//frame.Errors(frame.UserError, e.Error())
		return q
	}
	q.result = r
	l := len(fs)
	if l > 0 {
		for i := 0; i < l; i++ {
			parseResponseData(r, fs[i])
		}
		r.Close()
	}
	return q
}

func (m Model) Count() yorm.CountResult {
	m.s.GetStatement().Fields = `count(*) as dn_count`
	sql := query(m.s.GetStatement())
	c := cfr{}
	if m.s.GetStatement().SQL {
		c.sql = sql
	}
	if m.s.GetStatement().Exec {
		return c
	}
	r, e := m.s.Db().Query(sql)
	if e != nil {
		c.err = e
		return c
	}
	res := Count{}
	parseResponseData(r, &res)
	c.num = res.Count
	return c
}

func (m Model) Update(i interface{}) yorm.ExecResult {
	sql := update(m.s.GetStatement(), i)
	u := er{}
	if m.s.GetStatement().SQL {
		u.sql = sql
	}
	if m.s.GetStatement().Exec {
		return u
	}
	id, num, err := m.s.Db().Execute(sql)
	u.id = id
	u.num = num
	u.err = err
	return u
}

func (m Model) Insert(i interface{}) yorm.ExecResult {
	sql := exec(m.s.GetStatement(), i)
	e := er{}
	if m.s.GetStatement().SQL {
		e.sql = sql
	}
	if m.s.GetStatement().Exec {
		return e
	}
	id, num, err := m.s.Db().Execute(sql)
	e.id = id
	e.num = num
	e.err = err
	return e
}

func (m Model) Value(column string) yorm.ValueResult {
	if !strings.Contains(m.s.GetStatement().Fields, column) {
		m.s.Field(column)
	}
	sql := query(m.s.GetStatement())
	v := vfr{}
	if m.s.GetStatement().SQL {
		v.sql = sql
	}
	if m.s.GetStatement().Exec {
		return v
	}
	r, e := m.s.Db().Query(sql)
	if e != nil {
		v.err = e
		return v
	}
	raw := sql2.RawBytes{}
	r.Next()
	r.Scan(&raw)
	if len(raw) == 0 {
		v.nil = true
	} else {
		v.value = string(raw)
	}
	return v
}

func (m Model) Sum(column string) yorm.SumResult {
	m.s.GetStatement().Fields = `sum(` + column + `) as dn_sum`
	sql := query(m.s.GetStatement())
	s := sfr{}
	if m.s.GetStatement().SQL {
		s.sql = sql
	}
	if m.s.GetStatement().Exec {
		return s
	}
	r, e := m.s.Db().Query(sql)
	if e != nil {
		s.err = e
		return s
	}
	res := Sum{}
	parseResponseData(r, &res)
	s.num = res.Sum
	return s
}

func (m Model) Delete() yorm.ExecResult {
	sql := remove(m.s.GetStatement())
	e := er{}
	if m.s.GetStatement().SQL {
		e.sql = sql
	}
	if m.s.GetStatement().Exec {
		return e
	}
	id, num, err := m.s.Db().Execute(sql)
	e.id = id
	e.num = num
	e.err = err
	return e
}

// Transaction 开始事务
func (m Model) Transaction() (*yorm.Db, error) {
	return m.s.Db().Begin()
}

// Commit 提交事务
func (m Model) Commit() error {
	return m.s.Db().Commit()
}

// Rollback 回滚事务
func (m Model) Rollback() error {
	return m.s.Db().Rollback()
}
