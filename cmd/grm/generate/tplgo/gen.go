package tplgo

import (
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	ffmt "gopkg.in/ffmt.v1"
	grm "gopkg.in/grm.v1"
)

func Gen(limit, threads int, pkg, tag, base, out string) error {
	tpl := template.New("sql")
	tpl.Funcs(grm.Funcs)

	ff := []string{}
	filepath.Walk(base, func(path string, info os.FileInfo, err error) error {
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

	_, err := grm.ParseSqlFilesArgs(tpl, out, ff...)
	if err != nil {
		ffmt.Mark(err)
		return err
	}
	ms, err := ParseMethods(tpl.Templates())
	if err != nil {
		ffmt.Mark(err)
		return err
	}

	b := &TplData{
		Pkg:       pkg,
		By:        strings.Join(os.Args, " "),
		MaxLimit:  limit,
		MaxFork:   threads,
		FieldName: tag,
		Path:      base,
		Methods:   ms,
	}

	aaa := MakeTplData(b)
	aaa, err = format.Source(aaa)
	if err != nil {
		ffmt.Mark(err)
		return err
	}
	if out == "" {
		fmt.Println(string(aaa))
		return nil
	}

	or, _ := ioutil.ReadFile(out)
	if string(or) == string(aaa) {
		fmt.Println("Unchanged sql go file!")
		return nil
	}

	fmt.Println("Generate sql go file!")
	err = ioutil.WriteFile(out, aaa, 0666)
	if err != nil {
		ffmt.Mark(err)
		return err
	}
	return nil
}
