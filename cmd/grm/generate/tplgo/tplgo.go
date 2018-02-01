package tplgo

import (
	"bytes"
	"encoding/base64"
	"go/ast"
	"sort"
	"strings"
	"text/template"

	nomenclature "gopkg.in/go-grm/nomenclature.v1"
	grm "gopkg.in/grm.v1"
)

var _sql = `// Code generated by "{{.By}}"; DO NOT EDIT.

package {{.Pkg}}

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"text/template"
	"time"

	"gopkg.in/grm.v1"
	"gopkg.in/grm.v1/rows"
)

var (
{{if .Gray}}Path = "{{.Path}}"                                         // template file path {{end}}
	begin     = time.Now()                                             // start time
	MaxLimit  = {{.MaxLimit}}                                          // Max read rows limit
	FieldName = rows.MakeFieldName("{{.FieldName}}")                   // Field name
	MaxFork   = {{.MaxFork}}                                           // Max fork
	Template  = template.New("{{.Pkg}}")                               // template
	Stdout    = os.Stdout                                              // std out
	Log       = log.New(Stdout, "[grm] ", log.LstdFlags|log.Llongfile) // Print
	Println   = func(i ...interface{}) {
		Log.Output(3, fmt.Sprint(i...))
	} // logger
	GetDB = func() (*sql.DB, error) {
		db, err := grm.Get()
		return db, err
	} // db conn
)

func init() {
	Template.Funcs(grm.Funcs)
	{{if .Gray}}template.Must(grm.ParseSqlFiles(Template, Path)){{end}}
}

{{range .Methods}}

{{if eq .Type "Select"}}
// {{.Name}} {{.Comm}}
//line {{.Line}}
func {{.Name}}({{if .Req}}req *Req{{.Name}}, {{end}}dbs ...*sql.DB) (resp {{.Slice}}*Resp{{.Name}}, err error) {
	name := "{{.Name}}"
	
	temp := Template.Lookup(name)
	if temp == nil {
		var src []byte
		src, err = base64.StdEncoding.DecodeString("{{.Src}}")
		if err != nil {
			return
		}
		Template.New(name).Parse(string(src))
		temp = Template.Lookup(name)
	}
	
	var sqlStr string
	sqlStr, err = grm.Execute{{if .DDL}}DDL{{end}}(temp, {{if .Req}}req{{else}}nil{{end}})
	if err != nil {
		return 
	}
	
	if Println != nil {
		Println(sqlStr)
	}

	var db *sql.DB
	if len(dbs) != 0 {
		db = dbs[0]
	}

	if db == nil {
		db, err = GetDB()
		if err != nil {
			return
		}
	}

	_, err = grm.Query(db, sqlStr, {{if .Req}}req{{else}}nil{{end}}, &resp, MaxLimit, FieldName, MaxFork)
	return
}

{{if .Req}}
// Req{{.Name}} ...
//line {{.Line}}
type Req{{.Name}} struct { {{range .Req}}
	{{.Name}} {{.Type}} {{.Tags}} // {{.Comm}}{{end}}
}
{{end}}
// Resp{{.Name}} ...
//line {{.Line}}
type Resp{{.Name}} struct { {{range .Resp}}
	{{.Name}} {{.Type}} {{.Tags}} // {{.Comm}}{{end}}
}

{{if .Count}}
// {{.Name}}Count {{.Comm}}
//line {{.Line}}
func {{.Name}}Count({{if .ReqCount}}req *Req{{.Name}}Count, {{end}}dbs ...*sql.DB) (resp *Resp{{.Name}}Count, err error) {
	name := "{{.Name}}"
	
	temp := Template.Lookup(name)
	if temp == nil {
		var src []byte
		src, err = base64.StdEncoding.DecodeString("{{.Src}}")
		if err != nil {
			return
		}
		Template.New(name).Parse(string(src))
		temp = Template.Lookup(name)
	}
	
	var sqlStr string
	sqlStr, err = grm.ExecuteCount{{if .DDL}}DDL{{end}}(temp, {{if .ReqCount}}req{{else}}nil{{end}})
	if err != nil {
		return 
	}
	
	if Println != nil {
		Println(sqlStr)
	}

	var db *sql.DB
	if len(dbs) != 0 {
		db = dbs[0]
	}

	if db == nil {
		db, err = GetDB()
		if err != nil {
			return
		}
	}

	_, err = grm.Query(db, sqlStr, {{if .ReqCount}}req{{else}}nil{{end}}, &resp, MaxLimit, FieldName, MaxFork)
	return
}

{{if .ReqCount}}
// Req{{.Name}}Count ...
//line {{.Line}}
type Req{{.Name}}Count struct { {{range .ReqCount}}
	{{.Name}} {{.Type}} {{.Tags}} // {{.Comm}}{{end}}
}
{{end}}
// Resp{{.Name}}Count ...
//line {{.Line}}
type Resp{{.Name}}Count struct {
	Count int ` + "`" + `sql:"count"` + "`" + `
}
{{end}}

{{else if eq .Type "Update"}}
// {{.Name}} {{.Comm}}
//line {{.Line}}
func {{.Name}}({{if .Req}}req *Req{{.Name}}, {{end}}dbs ...*sql.DB) (count int,err error) {
	name := "{{.Name}}"
	
	temp := Template.Lookup(name)
	if temp == nil {
		var src []byte
		src, err = base64.StdEncoding.DecodeString("{{.Src}}")
		if err != nil {
			return
		}
		Template.New(name).Parse(string(src))
		temp = Template.Lookup(name)
	}
	
	var sqlStr string
	sqlStr, err = grm.Execute{{if .DDL}}DDL{{end}}(temp, {{if .Req}}req{{else}}nil{{end}})
	if err != nil {
		return 
	}
	
	if Println != nil {
		Println(sqlStr)
	}
	
	var db *sql.DB
	if len(dbs) != 0 {
		db = dbs[0]
	}

	if db == nil {
		db, err = GetDB()
		if err != nil {
			return
		}
	}

	return grm.ExecRowsAffected(db, sqlStr, {{if .Req}}req{{else}}nil{{end}})
}
{{if .Req}}
// Req{{.Name}} ...
//line {{.Line}}
type Req{{.Name}} struct { {{range .Req}}
	{{.Name}} {{.Type}} {{.Tags}} // {{.Comm}}{{end}}
}
{{end}}
{{else if eq .Type "Delete"}}
// {{.Name}} {{.Comm}}
//line {{.Line}}
func {{.Name}}({{if .Req}}req *Req{{.Name}}, {{end}}dbs ...*sql.DB) (count int,err error) {
	name := "{{.Name}}"
	
	temp := Template.Lookup(name)
	if temp == nil {
		var src []byte
		src, err = base64.StdEncoding.DecodeString("{{.Src}}")
		if err != nil {
			return
		}
		Template.New(name).Parse(string(src))
		temp = Template.Lookup(name)
	}
	
	var sqlStr string
	sqlStr, err = grm.Execute{{if .DDL}}DDL{{end}}(temp, {{if .Req}}req{{else}}nil{{end}})
	if err != nil {
		return 
	}
	
	if Println != nil {
		Println(sqlStr)
	}
	
	var db *sql.DB
	if len(dbs) != 0 {
		db = dbs[0]
	}

	if db == nil {
		db, err = GetDB()
		if err != nil {
			return
		}
	}

	return grm.ExecRowsAffected(db, sqlStr, {{if .Req}}req{{else}}nil{{end}})
}
{{if .Req}}
// Req{{.Name}} ...
//line {{.Line}}
type Req{{.Name}} struct { {{range .Req}}
	{{.Name}} {{.Type}} {{.Tags}} // {{.Comm}}{{end}}
}
{{end}}
{{else if eq .Type "Insert"}}
// {{.Name}} {{.Comm}}
//line {{.Line}}
func {{.Name}}({{if .Req}}req *Req{{.Name}}, {{end}}dbs ...*sql.DB) (count int,err error) {
	name := "{{.Name}}"
	
	temp := Template.Lookup(name)
	if temp == nil {
		var src []byte
		src, err = base64.StdEncoding.DecodeString("{{.Src}}")
		if err != nil {
			return
		}
		Template.New(name).Parse(string(src))
		temp = Template.Lookup(name)
	}
	
	var sqlStr string
	sqlStr, err = grm.Execute{{if .DDL}}DDL{{end}}(temp, {{if .Req}}req{{else}}nil{{end}})
	if err != nil {
		return 
	}
	
	if Println != nil {
		Println(sqlStr)
	}
	
	var db *sql.DB
	if len(dbs) != 0 {
		db = dbs[0]
	}

	if db == nil {
		db, err = GetDB()
		if err != nil {
			return
		}
	}
	return grm.ExecLastInsertId(db, sqlStr, {{if .Req}}req{{else}}nil{{end}})
}
{{if .Req}}
// Req{{.Name}} ...
//line {{.Line}}
type Req{{.Name}} struct { {{range .Req}}
	{{.Name}} {{.Type}} {{.Tags}} // {{.Comm}}{{end}}
}
{{end}}
{{else if eq .Type "Exec"}}
// {{.Name}} {{.Comm}}
//line {{.Line}}
func {{.Name}}({{if .Req}}req *Req{{.Name}}, {{end}}dbs ...*sql.DB) (err error) {
	name := "{{.Name}}"
	
	temp := Template.Lookup(name)
	if temp == nil {
		var src []byte
		src, err = base64.StdEncoding.DecodeString("{{.Src}}")
		if err != nil {
			return
		}
		Template.New(name).Parse(string(src))
		temp = Template.Lookup(name)
	}
	
	var sqlStr string
	sqlStr, err = grm.Execute{{if .DDL}}DDL{{end}}(temp, {{if .Req}}req{{else}}nil{{end}})
	if err != nil {
		return 
	}
	
	if Println != nil {
		Println(sqlStr)
	}
	
	var db *sql.DB
	if len(dbs) != 0 {
		db = dbs[0]
	}

	if db == nil {
		db, err = GetDB()
		if err != nil {
			return
		}
	}

	_, err = grm.Exec(db, sqlStr, {{if .Req}}req{{else}}nil{{end}})
	return err
}
{{if .Req}}
// Req{{.Name}} ...
//line {{.Line}}
type Req{{.Name}} struct { {{range .Req}}
	{{.Name}} {{.Type}} {{.Tags}} // {{.Comm}}{{end}}
}
{{end}}
{{end}}
{{end}}
`

var tpl, _ = template.New("").Parse(_sql)

func MakeTplData(data *TplData) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	err := tpl.Execute(buf, data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

type TplData struct {
	Gray      bool      // 启用灰度更新
	Pkg       string    // 包名
	By        string    // gen 注释
	MaxLimit  int       // 最大读取行数限制
	MaxFork   int       // 读取最大线程限制
	Path      string    // sql文件引用路径
	FieldName string    // 标签名
	Methods   []*Method // 方法
}

type Method struct {
	Src      string       // sql 源文件代码
	Line     string       // 第几行
	Name     string       // 定义的函数名
	Slice    string       // 类型前缀
	Type     string       // 操作类型
	Comm     string       // 注释
	Req      []*Parameter // 请求的参数
	ReqCount []*Parameter // 如果是获取count 请求的参数
	Resp     []*Parameter // 返回的参数
	Count    bool         // 是获取 count
	DDL      bool         // 是 DDL
}

type Parameter struct {
	Name string // 名字
	Type string // 类型
	Tags string // 标签
	Comm string // 注释
}

func ParseMethods(t []*template.Template) ([]*Method, error) {
	b := []*Method{}
	for _, v := range t {
		name := v.Name()
		parseName := v.ParseName
		if parseName == name || !ast.IsExported(name) {
			continue
		}
		con := v.Tree.Root.String()
		ss := strings.Split(con, "\n")

		ss0 := []string{}
		for _, v := range ss {
			if strings.HasPrefix(v, "--") {
				if len(v) == 2 {
					break
				}
				ss0 = append(ss0, v[2:])
			}
		}

		l, _ := v.ErrorContext(v.Tree.Root)

		if ci := strings.LastIndex(l, ":"); ci >= 0 {
			l = l[:ci]
		}
		m := &Method{
			Src:  base64.StdEncoding.EncodeToString([]byte(v.Root.String())),
			Line: l,
			Name: v.Tree.Name,
		}

		dd, err := grm.ReadAtLine(bytes.NewBufferString(strings.Join(ss0, "\n")))
		if err != nil {
			return nil, err
		}

		for _, v := range dd {
			if len(v) < 1 {
				continue
			}
			switch v[0] {
			case "@Count":
				m.Count = true
			case "@DDL":
				m.DDL = true
			case "@Type":
				if len(v) >= 2 {
					m.Type = v[1]
				}
				if len(v) >= 3 {
					m.Slice = v[2]
				}
			case "@Comm":
				if len(v) >= 2 {
					m.Comm += v[1]
					m.Comm += " "
				}

			case "@Req":
				if len(v) >= 4 {
					r := NewParameter(v)
					m.Req = append(m.Req, r)
					if m.Count {
						m.ReqCount = append(m.ReqCount, r)
					}
				}
			case "@Resp":
				if len(v) >= 4 {
					m.Resp = append(m.Resp, NewParameter(v))
				}
			}
		}

		b = append(b, m)
	}

	return b, nil
}

func NewParameter(v []string) *Parameter {
	ts := []string{}
	com := ""
	if len(v) >= 3 {
		ts = v[3 : len(v)-1]
		com = v[len(v)-1]
	}
	b := true

	t := "sql:"
	for _, v := range ts {
		if strings.Index(v, t) == 0 {
			b = false
			break
		}
	}

	if b {
		ts = append(ts, t+`"`+nomenclature.Hump2Snake(v[1])+`"`)
	}

	sort.Strings(ts)
	return &Parameter{
		Name: v[1],
		Type: v[2],
		Tags: "`" + strings.Join(ts, " ") + "`",
		Comm: com,
	}
}
