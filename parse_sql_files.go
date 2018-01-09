package grm

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

func ParseSqlFiles(t *template.Template, commit int, paths ...string) (*template.Template, error) {
	ext := ".sql"
	filenames := []string{}
	for _, path := range paths {
		filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			if filepath.Ext(path) != ext {
				return nil
			}
			filenames = append(filenames, path[:len(path)-len(ext)])
			return nil
		})
	}
	if t == nil {
		t = template.New("_")
	}

	mc := regexp.MustCompile("--.*\n")
	for _, filename := range filenames {
		b, err := ioutil.ReadFile(filename + ext)
		if err != nil {
			return nil, err
		}
		if commit == 1 {
			b = mc.ReplaceAll(b, []byte("\n"))
		}
		s := string(b)

		filename = strings.Replace(filename, `\`, `/`, -1)
		tmpl := t.New(filename)
		_, err = tmpl.Parse(s)
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}
