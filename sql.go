package grm

import (
	"bytes"
	"context"
	"database/sql"
	"io"
	"reflect"
	"sync"

	"gopkg.in/go-grm/rows.v1"
)

var pool = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(nil)
	},
}

type BaseData struct {
	Name string
	Data interface{}
}

type TemplateExecute interface {
	Execute(wr io.Writer, data interface{}) error
}

type DBQuery interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

func Query(db DBQuery, tpl TemplateExecute, req, resp interface{},
	limit int, fn func(reflect.StructField) string, f int) (int, error) {
	buf := pool.Get().(*bytes.Buffer)
	buf.Reset()

	err := tpl.Execute(buf, req)
	if err != nil {
		pool.Put(buf)
		return 0, err
	}
	sqlStr := buf.String()
	pool.Put(buf)

	row, err := db.Query(sqlStr)
	if err != nil {
		return 0, err
	}
	defer row.Close()
	i, err := rows.RowsScan(row, resp, limit, fn, f)
	if err != nil {
		return 0, err
	}
	return i, nil
}

type DBExec interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

func Exec(db DBExec, tpl TemplateExecute, req interface{}) (sql.Result, error) {
	buf := pool.Get().(*bytes.Buffer)
	buf.Reset()

	err := tpl.Execute(buf, req)
	if err != nil {
		pool.Put(buf)
		return nil, err
	}
	sqlStr := buf.String()
	pool.Put(buf)
	return db.Exec(sqlStr)
}

func ExecRowsAffected(db DBExec, tpl TemplateExecute, req interface{}) (int, error) {
	resu, err := Exec(db, tpl, req)
	if err != nil {
		return 0, err
	}
	aff, err := resu.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(aff), nil
}

func ExecLastInsertId(db DBExec, tpl TemplateExecute, req interface{}) (int, error) {
	resu, err := Exec(db, tpl, req)
	if err != nil {
		return 0, err
	}
	aff, err := resu.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(aff), nil
}
