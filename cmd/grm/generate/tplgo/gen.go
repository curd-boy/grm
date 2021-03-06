package tplgo

import (
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	grm "gopkg.in/grm.v1"
)

func Gen(limit, threads int, pkg, tag, base, out string, grayscale bool) error {
	tpl := template.New("sql")
	tpl.Funcs(grm.Funcs)

	ff := []string{}
	err := filepath.Walk(base, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".sql" {
			return nil
		}
		ff = append(ff, path)
		return nil
	})
	if err != nil {
		return err
	}

	// 解析文件
	_, err = grm.ParseSqlFilesArgs(tpl, out, ff...)
	if err != nil {
		return err
	}

	tt := tpl.Templates()
	// 排序模板
	sort.Slice(tt, func(i, j int) bool {
		return tt[i].Name() <= tt[j].Name()
	})

	// 解析模板注解
	ms, err := ParseMethods(tt)
	if err != nil {
		return err
	}

	// 拼凑模板数据
	b := &TplData{
		Gray:      grayscale,
		Pkg:       pkg,
		By:        strings.Join(os.Args, " "),
		MaxLimit:  limit,
		MaxFork:   threads,
		FieldName: tag,
		Path:      base,
		Methods:   ms,
	}

	// 填充数据
	src, err := MakeTplData(b)
	if err != nil {
		return err
	}

	// 代码格式化
	src, err = format.Source(src)
	if err != nil {
		return err
	}

	// 是否直接输出
	if out == "" {
		fmt.Println(string(src))
		return nil
	}

	out = path.Join(out)
	// 比较文件
	or, _ := ioutil.ReadFile(out)
	if string(or) == string(src) {
		fmt.Println("[grm] Unchanged " + out)
		return nil
	}

	// 如果不同则修改
	fmt.Println("[grm] Generate " + out)
	err = ioutil.WriteFile(out, src, 0666)
	if err != nil {
		return err
	}
	return nil
}
