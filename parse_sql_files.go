package grm

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

func ParseSqlFilesArgs(t *template.Template, out string, paths ...string) (*template.Template, error) {
	if out != "" {
		out = filepath.Dir(filepath.Clean(out))
	}
	return parseSqlFiles(t, 0, out, paths)
}

func ParseSqlFiles(t *template.Template, paths ...string) (*template.Template, error) {
	return parseSqlFiles(t, 1, "", paths)
}

func parseSqlFiles(t *template.Template, commit int, out string, paths []string) (*template.Template, error) {
	ext := ".sql"
	filenames := []string{}
	for _, path := range paths {
		filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			if filepath.Ext(path) != ext {
				return nil
			}
			filenames = append(filenames, path)
			return nil
		})
	}
	if t == nil {
		t = template.New("_")
	}

	mc := regexp.MustCompile("--.*\n")
	for _, filename := range filenames {
		b, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}
		if commit == 1 {
			b = mc.ReplaceAll(b, []byte("\n"))
		}
		s := string(b)

		if out != "" {
			filename = filepath.Clean(filename)
			filename0, err := filepath.Rel(out, filename)
			if err == nil {
				filename = filename0
			}
		}

		filename = strings.Replace(filename, `\`, `/`, -1)

		tmpl := t.New(filename)
		_, err = tmpl.Parse(s)
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}
