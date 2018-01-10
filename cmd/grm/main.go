package main

import (
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"gopkg.in/ffmt.v1"
	"gopkg.in/grm.v1"
	"gopkg.in/grm.v1/cli"
	"gopkg.in/grm.v1/cmd/grm/logo"
)

func main() {

	app := cli.App{}

	app.Name = "grm"

	app.Usage = "Grm is a tool like mybatis for the Go programming language."

	app.Version = version

	app.Authors = []*cli.Author{
		{
			Name:  "wzshiming",
			Email: "wzshiming@foxmail.com",
		},
	}

	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:    "hidden-logo",
			Aliases: []string{"hl"},
			Usage:   "hidden logo",
			Value:   false,
		},
	}

	app.Before = func(c *cli.Context) error {
		if !c.Bool("hidden-logo") {
			logo.PrintLogo("V" + version)
		}
		return nil
	}

	app.Commands = Commands()

	app.Run(os.Args)
}

func Commands() []*cli.Command {
	r := []*cli.Command{
		{
			Name:        "generate",
			Aliases:     []string{"gen", "g"},
			Usage:       "Generate commands",
			Subcommands: SubcommandsGenerate(),
		},
	}
	return r
}

func SubcommandsGenerate() []*cli.Command {
	return []*cli.Command{
		{
			Name:    "go",
			Aliases: []string{"g"},
			Usage:   "Generate go file",
			Flags: []cli.Flag{
				&cli.IntFlag{
					Name:    "limit",
					Aliases: []string{"l"},
					Usage:   "Maximum limit rows",
					Value:   10000,
				},
				&cli.IntFlag{
					Name:    "threads",
					Aliases: []string{"t"},
					Usage:   "A sql threads size",
					Value:   3,
				},
				&cli.StringFlag{
					Name:  "tag",
					Usage: "Tag name",
					Value: "sql",
				},
				&cli.StringFlag{
					Name:    "filepath",
					Aliases: []string{"f"},
					Usage:   "Sql file path",
					Value:   "./",
				},
				&cli.StringFlag{
					Name:    "package",
					Aliases: []string{"p"},
					Usage:   "Package name",
					Value:   "sql",
				},
				&cli.StringFlag{
					Name:    "out",
					Aliases: []string{"o"},
					Usage:   "out file",
					Value:   "",
				},
			},
			Action: func(c *cli.Context) error {
				limit := c.Int("limit")
				threads := c.Int("threads")
				pkg := c.String("package")
				tag := c.String("tag")
				path := c.String("filepath")
				out := c.String("out")
				run(limit, threads, pkg, tag, path, out)
				return nil
			},
		},
	}
}

func run(limit, threads int, pkg, tag, base, out string) {
	tpl := template.New("sql")
	tpl.Funcs(grm.Funcs)

	ff := []string{}
	err := filepath.Walk(base, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
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
		return
	}

	_, err = grm.ParseSqlFiles(tpl, 0, ff...)
	if err != nil {
		ffmt.Mark(err)
		return
	}
	ms, err := ParseMethods(tpl.Templates())
	if err != nil {
		ffmt.Mark(err)
		return
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
		return
	}
	if out == "" {
		fmt.Println(string(aaa))
		return
	}

	or, _ := ioutil.ReadFile(out)
	if string(or) == string(aaa) {
		fmt.Println("Unchanged sql go file!")
		return
	}

	fmt.Println("Generate sql go file!")
	err = ioutil.WriteFile(out, aaa, 0666)
	if err != nil {
		ffmt.Mark(err)
		return
	}
}
