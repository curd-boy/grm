package main

import (
	"os"

	"gopkg.in/grm.v1/cli"
	sqlfmt "gopkg.in/grm.v1/cmd/grm/format"
	"gopkg.in/grm.v1/cmd/grm/generate/tplgo"
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
		{
			Name:    "format",
			Aliases: []string{"fmt", "f"},
			Usage:   "Format sql file",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "filepath",
					Aliases: []string{"f"},
					Usage:   "Sql file path",
					Value:   "./",
				},
			},
			Action: func(c *cli.Context) error {

				path := c.String("filepath")
				return sqlfmt.FormatDir(path)
			},
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
				tplgo.Gen(limit, threads, pkg, tag, path, out)
				return nil
			},
		},
	}
}
