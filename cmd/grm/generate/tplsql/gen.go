package tplsql

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"
	"text/template"

	ffmt "gopkg.in/ffmt.v1"
	nomenclature "gopkg.in/go-grm/nomenclature.v1"
	grm "gopkg.in/grm.v1"
	sqlfmt "gopkg.in/grm.v1/cmd/grm/format"
)

//go:generate grm --hl fmt -f ./sql/
//go:generate grm --hl gen go -p tplsql -f ./sql/ -o ./sql.go

func Gen(conn string, out string) error {
	ur, err := url.Parse(conn)
	if err != nil {
		return err
	}

	if ur.Opaque != "" {
		return errors.New("error db conn :" + conn)
	}

	// 打开 数据库连接
	_, err = grm.Register(conn)
	if err != nil {
		return err
	}

	// 获取当前库库
	s, err := GetSchema(nil)
	if err != nil {
		return err
	}

	// 获取表
	t, err := GetTable(nil, &ReqGetTable{
		TableSchema: s.TableSchema,
	})
	if err != nil {
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
		resp, err := GetCreateTable(nil, &ReqGetCreateTable{
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
			Name: "CreateTable" + nomenclature.Snake2Hump(v.TableName),
			Sql:  resp.SqlCreateTable,
		})

		col, err := GetColumn(nil, &ReqGetColumn{
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
	src := sqlfmt.Format(buf.Bytes())

	if out == "" {
		fmt.Println(string(src))
		return nil
	}

	or, _ := ioutil.ReadFile(out)
	if string(or) == string(src) {
		fmt.Println("Unchanged sql file!")
		return nil
	}

	fmt.Println("Generate sql  file!")
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
		k0 = append(k0, "`"+v.ColumnName+"`")
		k1 = append(k1, ":"+nomenclature.Snake2Hump(v.ColumnName))
		d = append(d, &ParameterTplData{
			Method:    "@Req",
			Name:      nomenclature.Snake2Hump(v.ColumnName),
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
		Name:      "Inster" + nomenclature.Snake2Hump(table),
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
				Name:      nomenclature.Snake2Hump(oriName),
				OriAsName: oriName,
				OriName:   v.ColumnName,
				Type:      GetDataType2GoType(v.DataType),
				Comm:      v.ColumnComment,
			})
		}

		k0 = append(k0, oriAs)
		resp = append(resp, &ParameterTplData{
			Method:    "@Resp",
			Name:      nomenclature.Snake2Hump(oriName),
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
		Name:      "SelectFirst" + nomenclature.Snake2Hump(table),
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
		newName := nomenclature.Snake2Hump(oriName)

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
		Name:      "UpdateFirst" + nomenclature.Snake2Hump(table),
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

		ptd := ParameterTplData{
			Method:    "@Resp",
			Name:      nomenclature.Snake2Hump(oriName),
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
		Type:      "Select",
		Comm:      "Select limit offset " + table,
		Name:      "SelectAll" + nomenclature.Snake2Hump(table),
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

		ptd := ParameterTplData{
			Method:    "@Req",
			Name:      nomenclature.Snake2Hump(oriName),
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
		Name:      "DeleteFirst" + nomenclature.Snake2Hump(table),
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
