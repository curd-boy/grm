package tplgo

import (
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	ffmt "gopkg.in/ffmt.v1"
	grm "gopkg.in/grm.v1"
)

func Gen(limit, threads int, pkg, tag, base, out string) error {
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
		ffmt.Mark(err)
		return err
	}

	_, err = grm.ParseSqlFilesArgs(tpl, out, ff...)
	if err != nil {
		ffmt.Mark(err)
		return err
	}

	tt := tpl.Templates()

	sort.Slice(tt, func(i, j int) bool {
		return tt[i].Name() <= tt[j].Name()
	})

	ms, err := ParseMethods(tt)
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

	src, err := MakeTplData(b)
	if err != nil {
		ffmt.Mark(err)
		return err
	}
	src, err = format.Source(src)
	if err != nil {
		ffmt.Mark(err)
		return err
	}
	if out == "" {
		fmt.Println(string(src))
		return nil
	}

	or, _ := ioutil.ReadFile(out)
	if string(or) == string(src) {
		fmt.Println("Unchanged sql go file!")
		return nil
	}

	fmt.Println("Generate sql go file!")
	err = ioutil.WriteFile(out, src, 0666)
	if err != nil {
		ffmt.Mark(err)
		return err
	}
	return nil
}
