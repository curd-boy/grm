package tplsql

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"path"
	"strings"
	"text/template"

	ffmt "gopkg.in/ffmt.v1"
	grm "gopkg.in/grm.v1"
	sqlfmt "gopkg.in/grm.v1/cmd/grm/format"
	namecase "gopkg.in/wzshiming/namecase.v2"
)

//go:generate grm --hl fmt -f ./sql/
//go:generate grm --hl gen go -p tplsql -f ./sql/ -o ./sql.go

func Gen(conn string, out string) error {
	ur, err := url.Parse(conn)
	if err != nil {
		ffmt.Mark(err)
		return err
	}

	if ur.Opaque != "" {
		return errors.New("error db conn :" + conn)
	}

	// 打开 数据库连接
	_, err = grm.Register(conn)
	if err != nil {
		ffmt.Mark(err)
		return err
	}

	// 禁止输出
	Println = nil

	// 获取当前库库
	s, err := GetSchema()
	if err != nil {
		ffmt.Mark(err)
		return err
	}

	// 获取表
	t, err := GetTable(&ReqGetTable{
		TableSchema: s.TableSchema,
	})
	if err != nil {
		ffmt.Mark(err)
		return err
	}

	// 初始化模板
	temp0 := template.New("tplsql")
	temp, err := temp0.Parse(_sql)
	if err != nil {
		ffmt.Mark(err)
		return err
	}
	buf := bytes.NewBuffer(nil)

	ttd := []*DefinesTplData{}
	for _, v := range t {
		resp, err := GetCreateTable(&ReqGetCreateTable{
			TableSchema: s.TableSchema,
			TableName:   v.TableName,
		})
		if err != nil {
			ffmt.Mark(err)
			return err
		}

		// 生成创建表的sql
		ttd = append(ttd, &DefinesTplData{
			Type: "Exec",
			Comm: "Create table " + v.TableName,
			Name: "CreateTable" + namecase.ToUpperHump(v.TableName),
			Sql:  resp.SqlCreateTable,
		})

		col, err := GetColumn(&ReqGetColumn{
			TableSchema: s.TableSchema,
			TableName:   v.TableName,
		})
		if err != nil {
			ffmt.Mark(err)
			return err
		}

		ttd = append(ttd, MakeInsterInfo(v.TableName, col))
		ttd = append(ttd, MakeDeleteFirst(v.TableName, col))
		ttd = append(ttd, MakeSelectFirst(v.TableName, col))
		ttd = append(ttd, MakeSelectAll(v.TableName, col))
		ttd = append(ttd, MakeUpdateFirst(v.TableName, col))
	}
	buf.Reset()
	err = temp.Execute(buf, ttd)
	if err != nil {
		ffmt.Mark(err)
		return err
	}
	//格式化
	src := sqlfmt.Format(buf.Bytes())

	// 直接输出
	if out == "" {
		fmt.Println(string(src))
		return nil
	}

	// 比较稳健
	out = path.Join(out)
	or, _ := ioutil.ReadFile(out)
	if string(or) == string(src) {
		fmt.Println("[grm] Unchanged " + out)
		return nil
	}

	fmt.Println("[grm] Generate " + out)
	err = ioutil.WriteFile(out, src, 0666)
	if err != nil {
		ffmt.Mark(err)
		return err
	}
	return nil
}

func MakeInsterInfo(table string, col []*RespGetColumn) *DefinesTplData {
	k0 := []string{}
	k1 := []string{}
	d := []*ParameterTplData{}
	for _, v := range col {
		if v.ColumnKey == "PRI" {
			continue
		}

		v.ColumnComment = strings.Replace(v.ColumnComment, "\n", " ", -1)
		k0 = append(k0, "`"+v.ColumnName+"`")
		k1 = append(k1, ":"+namecase.ToUpperHump(v.ColumnName))
		d = append(d, &ParameterTplData{
			Method:    "@Req",
			Name:      namecase.ToUpperHump(v.ColumnName),
			OriAsName: v.ColumnName,
			OriName:   v.ColumnName,
			Type:      GetDataType2GoType(v.DataType),
			Comm:      v.ColumnComment,
		})

	}
	t := fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s);`, "`"+table+"`", strings.Join(k0, ", "), strings.Join(k1, ", "))

	return &DefinesTplData{
		Type:      "Insert",
		Comm:      "Inster info " + table,
		Name:      "Inster" + namecase.ToUpperHump(table),
		Sql:       t,
		Parameter: d,
	}
}

func MakeSelectFirst(table string, col []*RespGetColumn) *DefinesTplData {
	k0 := []string{}
	resp := []*ParameterTplData{}
	req := []*ParameterTplData{}
	for _, v := range col {
		oriName := v.ColumnName
		oriAs := "`" + v.ColumnName + "`"
		if oriName == "id" {
			oriName = table + "_id"
			oriAs = "`id` AS `" + oriName + "`"
		}

		if v.ColumnKey == "PRI" {
			req = append(req, &ParameterTplData{
				Method:    "@Req",
				Name:      namecase.ToUpperHump(oriName),
				OriAsName: oriName,
				OriName:   v.ColumnName,
				Type:      GetDataType2GoType(v.DataType),
				Comm:      v.ColumnComment,
			})
		}

		v.ColumnComment = strings.Replace(v.ColumnComment, "\n", " ", -1)
		k0 = append(k0, oriAs)
		resp = append(resp, &ParameterTplData{
			Method:    "@Resp",
			Name:      namecase.ToUpperHump(oriName),
			OriAsName: oriName,
			OriName:   v.ColumnName,
			Type:      GetDataType2GoType(v.DataType),
			Comm:      v.ColumnComment,
		})

	}
	where := ""
	if len(req) != 0 {
		where = fmt.Sprintf(" WHERE `%s` = :%s", req[0].OriAsName, req[0].Name)
	}
	t := fmt.Sprintf(`SELECT %s FROM %s%s LIMIT 1;`, strings.Join(k0, ", "), "`"+table+"`", where)

	return &DefinesTplData{
		Type:      "Select",
		Comm:      "Select first " + table,
		Name:      "SelectFirst" + namecase.ToUpperHump(table),
		Sql:       t,
		Parameter: append(req, resp...),
	}
}

func MakeUpdateFirst(table string, col []*RespGetColumn) *DefinesTplData {
	k0 := []string{}
	req1 := []*ParameterTplData{}
	req := []*ParameterTplData{}
	for _, v := range col {
		oriName := v.ColumnName
		oriAsName := "`" + v.ColumnName + "`"
		if oriName == "id" {
			oriName = table + "_id"
		}
		newName := namecase.ToUpperHump(oriName)

		if v.ColumnKey == "PRI" {
			req = append(req, &ParameterTplData{
				Method:    "@Req",
				Name:      newName,
				OriAsName: oriName,
				OriName:   v.ColumnName,
				Type:      GetDataType2GoType(v.DataType),
				Comm:      v.ColumnComment,
			})
			continue
		}

		v.ColumnComment = strings.Replace(v.ColumnComment, "\n", " ", -1)
		k0 = append(k0, oriAsName+" = :"+newName)
		req1 = append(req1, &ParameterTplData{
			Method:    "@Req",
			Name:      newName,
			OriAsName: oriName,
			OriName:   v.ColumnName,
			Type:      GetDataType2GoType(v.DataType),
			Comm:      v.ColumnComment,
		})

	}
	where := ""
	if len(req) != 0 {
		where = fmt.Sprintf(" WHERE `%s` = :%s", req[0].OriName, req[0].Name)
	}
	t := fmt.Sprintf(`UPDATE %s SET %s%s LIMIT 1;`, "`"+table+"`", strings.Join(k0, ", "), where)
	return &DefinesTplData{
		Type:      "Update",
		Comm:      "Update first " + table,
		Name:      "UpdateFirst" + namecase.ToUpperHump(table),
		Sql:       t,
		Parameter: append(req, req1...),
	}
}

func MakeSelectAll(table string, col []*RespGetColumn) *DefinesTplData {
	k0 := []string{}
	resp := []*ParameterTplData{}
	req := []*ParameterTplData{}
	req = append(req, &ParameterTplData{
		Method:    "@Req",
		Name:      "Offset",
		OriAsName: "offset",
		OriName:   "offset",
		Type:      "int",
		Comm:      "offset index",
	})
	req = append(req, &ParameterTplData{
		Method:    "@Req",
		Name:      "Limit",
		OriAsName: "limit",
		OriName:   "limit",
		Type:      "int",
		Comm:      "limit rows",
	})
	req = append(req, &ParameterTplData{
		Method: "@Count",
	})

	for _, v := range col {
		oriName := v.ColumnName
		oriAs := "`" + v.ColumnName + "`"
		if oriName == "id" {
			oriName = table + "_id"
			oriAs = "`id` AS `" + oriName + "`"
		}

		v.ColumnComment = strings.Replace(v.ColumnComment, "\n", " ", -1)
		ptd := ParameterTplData{
			Method:    "@Resp",
			Name:      namecase.ToUpperHump(oriName),
			OriAsName: oriName,
			OriName:   v.ColumnName,
			Type:      GetDataType2GoType(v.DataType),
			Comm:      v.ColumnComment,
		}

		k0 = append(k0, oriAs)
		resp = append(resp, &ptd)

	}

	t := fmt.Sprintf(`SELECT %s FROM %s LIMIT :Offset, :Limit;`, strings.Join(k0, ", "), "`"+table+"`")

	return &DefinesTplData{
		Type:      "Select []",
		Comm:      "Select limit offset " + table,
		Name:      "SelectAll" + namecase.ToUpperHump(table),
		Sql:       t,
		Parameter: append(req, resp...),
	}
}

func MakeDeleteFirst(table string, col []*RespGetColumn) *DefinesTplData {

	req := []*ParameterTplData{}
	for _, v := range col {
		oriName := v.ColumnName
		if oriName == "id" {
			oriName = table + "_id"
		}

		v.ColumnComment = strings.Replace(v.ColumnComment, "\n", " ", -1)
		ptd := ParameterTplData{
			Method:    "@Req",
			Name:      namecase.ToUpperHump(oriName),
			OriAsName: oriName,
			OriName:   v.ColumnName,
			Type:      GetDataType2GoType(v.DataType),
			Comm:      v.ColumnComment,
		}
		if v.ColumnKey == "PRI" {
			req = append(req, &ptd)
		}

	}
	where := ""
	if len(req) != 0 {
		where = fmt.Sprintf(" WHERE `%s` = :%s", req[0].OriName, req[0].Name)
	}
	t := fmt.Sprintf(`DELETE FROM %s%s LIMIT 1;`, "`"+table+"`", where)

	return &DefinesTplData{
		Type:      "Delete",
		Comm:      "Delete first " + table,
		Name:      "DeleteFirst" + namecase.ToUpperHump(table),
		Sql:       t,
		Parameter: req,
	}
}

func GetDataType2GoType(dt string) string {
	d, ok := mysqlTypes[dt]
	if ok {
		return d
	}
	return "string"
}

var mysqlTypes = map[string]string{
	"bool":      "bool",
	"timestamp": "time.Time",
	"date":      "time.Time",
	"datetime":  "time.Time",
	"time":      "time.Time",
	"int":       "int",
	"tinyint":   "int8",
	"smallint":  "int16",
	"integer":   "int32",
	"bigint":    "int64",
	"float":     "float32",
	"double":    "float64",
	"decimal":   "float64",
}
