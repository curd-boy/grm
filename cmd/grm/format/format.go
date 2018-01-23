package format

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"

	ffmt "gopkg.in/ffmt.v1"
	grm "gopkg.in/grm.v1"
)

func FormatDir(pa string) error {
	return filepath.Walk(pa, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".sql" {
			return nil
		}
		d, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		d0 := Format(d)
		if string(d0) == string(d) {
			return nil
		}
		err = ioutil.WriteFile(path, d0, 0666)
		if err != nil {
			return err
		}
		return nil
	})
}

func Format(src []byte) []byte {
	n := []byte("\n")
	ss := bytes.Split(src, n)

	dists := [][]byte{}
	d := bytes.NewBuffer(nil)
	for _, v := range ss {
		if bytes.Index(v, []byte("--")) != 0 {
			if d.Len() != 0 {
				ral, err := grm.ReadAtLine(d)
				if err != nil {
					ffmt.Mark(err)
				}
				ft := grm.WriterAtLine(ral)

				dists = append(dists, ft.Bytes())
				d.Reset()
			}
			dists = append(dists, v)
			continue
		}

		d.Write(v)
		d.WriteByte('\n')
	}
	return bytes.Join(dists, n)
}
