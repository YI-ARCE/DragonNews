package mysql

import (
	"yiarce/core/yorm"
)

func init() {
	yorm.Register("mysql", Transfer{})
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
	if m.s.GetStatement().Sql {
		q.sql = sql
	}
	if m.s.GetStatement().Exec {
		return &q
	}
	r, e := m.s.Db().Query(sql)
	if e != nil {
		q.err = e
		return &q
	}
	q.result = r
	if len(fs) > 0 {
		parseResponseData(r, fs[0])
	}
	return &q
}

func (m Model) Select(fs ...yorm.FormatStruct) yorm.QueryResults {
	sql := query(m.s.GetStatement())
	q := qr{}
	if m.s.GetStatement().Sql {
		q.sql = sql
	}
	if m.s.GetStatement().Exec {
		return &q
	}
	r, e := m.s.Db().Query(sql)
	if e != nil {
		q.err = e
		//frame.Errors(frame.UserError, e.Error(), nil)
		return &q
	}
	q.result = r
	if len(fs) > 0 {
		parseResponseData(r, fs[0])
	}
	return &q
}

func (m Model) Update(i interface{}) yorm.ExecResult {
	sql := update(m.s.GetStatement(), i)
	u := er{}
	if m.s.GetStatement().Sql {
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
	if m.s.GetStatement().Sql {
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

func (m Model) Delete() yorm.ExecResult {
	sql := remove(m.s.GetStatement())
	e := er{}
	if m.s.GetStatement().Sql {
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

func (m Model) Transaction() error {
	return nil
}

func (m Model) Commit() error {
	return nil
}

func (m Model) Rollback() error {
	return nil
}
