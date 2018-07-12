package grm

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"reflect"
	"sync"

	"github.com/wzshiming/rows"
	sqlparser "gopkg.in/go-grm/sqlparser.v1"
	"gopkg.in/go-grm/sqlparser.v1/dependency/querypb"
	"gopkg.in/go-grm/sqlparser.v1/dependency/sqltypes"
)

var pool = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(nil)
	},
}

type TemplateExecute interface {
	Execute(wr io.Writer, data interface{}) error
}

type DBQuery interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

func toBindVariable(args interface{}) (bindVariables map[string]*querypb.BindVariable, extras map[string]sqlparser.Encodable, err error) {
	v := reflect.ValueOf(args)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	bindVariables = map[string]*querypb.BindVariable{}
	extras = map[string]sqlparser.Encodable{}
	switch v.Kind() {
	case reflect.Struct:
		t := v.Type()
		j := t.NumField()
		for i := 0; i != j; i++ {
			tf := t.Field(i)
			vf := v.Field(i)
			bv, err := sqltypes.BuildBindVariable(vf.Interface())
			if err != nil {
				return bindVariables, extras, err
			}
			bindVariables[tf.Name] = bv
		}
	case reflect.Map:
		for _, k := range v.MapKeys() {
			vf := v.MapIndex(k)
			bv, err := sqltypes.BuildBindVariable(vf.Interface())
			if err != nil {
				return bindVariables, extras, err
			}
			bindVariables[fmt.Sprint(k.Interface())] = bv
		}
	}
	return bindVariables, extras, nil
}

func Execute(tpl TemplateExecute, req interface{}) (string, error) {
	return execute(tpl, req, false, false)
}

func ExecuteDDL(tpl TemplateExecute, req interface{}) (string, error) {
	return execute(tpl, req, false, true)
}

func ExecuteCount(tpl TemplateExecute, req interface{}) (string, error) {
	return execute(tpl, req, true, false)
}

func ExecuteDDLCount(tpl TemplateExecute, req interface{}) (string, error) {
	return execute(tpl, req, true, true)
}

// count(1) AS count
var cac = sqlparser.SelectExprs{
	&sqlparser.AliasedExpr{
		Expr: &sqlparser.FuncExpr{
			Name: sqlparser.NewColIdent("count"),
			Exprs: sqlparser.SelectExprs{
				&sqlparser.AliasedExpr{
					Expr: sqlparser.NewIntVal([]byte("1")),
				},
			},
		},
		As: sqlparser.NewColIdent("count"),
	},
}

func execute(tpl TemplateExecute, req interface{}, isCount, isDDL bool) (string, error) {
	if tpl == nil {
		return "", errors.New("Error Execute: Template is nil")
	}

	buf := pool.Get().(*bytes.Buffer)
	defer pool.Put(buf)

	buf.Reset()
	err := tpl.Execute(buf, req)
	if err != nil {
		return "", err
	}

	var stat sqlparser.Statement
	if isDDL {
		stat, err = sqlparser.ParseStrictDDL(buf.String())
		if err != nil {
			return "", err
		}
	} else {
		stat, err = sqlparser.Parse(buf.String())
		if err != nil {
			return "", err
		}
	}

	if isCount {
		switch st := stat.(type) {
		case *sqlparser.Select:
			st.Limit = nil
			st.OrderBy = nil
			st.SelectExprs = cac
		}
	}

	buf.Reset()
	tb := &sqlparser.TrackedBuffer{Buffer: buf}
	stat.Format(tb)

	// Not needed bind
	if !tb.HasBindVars() {
		return buf.String(), nil
	}

	if req == nil {
		return "", errors.New("Error Execute: No binding data")
	}

	b, e, err := toBindVariable(req)
	if err != nil {
		return "", err
	}

	pq := tb.ParsedQuery()

	d, err := pq.GenerateQuery(b, e)
	if err != nil {
		return "", err
	}

	return string(d), nil
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
