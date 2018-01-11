package grm

import (
	"bytes"
	"context"
	"database/sql"
	"io"
	"reflect"
	"sync"

	"gopkg.in/go-grm/rows.v1"
	"gopkg.in/go-grm/sqlparser.v1"
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

func Execute(tpl TemplateExecute, req interface{}) (string, error) {
	buf := pool.Get().(*bytes.Buffer)
	buf.Reset()
	defer pool.Put(buf)
	err := tpl.Execute(buf, req)
	if err != nil {
		return "", err
	}
	sqlStr := buf.String()

	stat, err := sqlparser.Parse(sqlStr)
	if err != nil {
		return "", err
	}

	sqlStr = sqlparser.String(stat)
	return sqlStr, nil
}

func Query(db DBQuery, sqlStr string, req, resp interface{},
	limit int, fn func(reflect.StructField) string, f int) (int, error) {
	row, err := db.Query(sqlStr)
	if err != nil {
		return 0, err
	}

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

func Exec(db DBExec, sqlStr string, req interface{}) (sql.Result, error) {
	return db.Exec(sqlStr)
}

func ExecRowsAffected(db DBExec, sqlStr string, req interface{}) (int, error) {
	resu, err := Exec(db, sqlStr, req)
	if err != nil {
		return 0, err
	}
	aff, err := resu.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(aff), nil
}

func ExecLastInsertId(db DBExec, sqlStr string, req interface{}) (int, error) {
	resu, err := Exec(db, sqlStr, req)
	if err != nil {
		return 0, err
	}
	aff, err := resu.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(aff), nil
}
