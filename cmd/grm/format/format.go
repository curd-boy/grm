package format

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	ffmt "gopkg.in/ffmt.v1"
	nomenclature "gopkg.in/go-grm/nomenclature.v1"
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

	dists0 := [][]byte{}
	d := bytes.NewBuffer(nil)
	for _, v := range ss {
		if bytes.Index(v, []byte("-- @")) != 0 {
			if d.Len() != 0 {
				ral, err := grm.ReadAtLine(d)
				if err != nil {
					ffmt.Mark(err)
				}
				for k0, v0 := range ral {
					if len(v0) > 2 {
						switch v0[1] {
						case "@Req", "@Resp":
							ral[k0][2] = nomenclature.Snake2Hump(ral[k0][2])
							if len(v0) < 4 {
								ral[k0] = append(ral[k0], "string")
							}

							if ral[k0][3] == "" {
								ral[k0][3] = "string"
							}

							if len(v0) < 5 {
								ral[k0] = append(ral[k0], ral[k0][2])
							}

							if ral[k0][4] == "" {
								ral[k0][4] = strings.Replace(nomenclature.Hump2Snake(ral[k0][2]), "_", " ", -1)
							}
						}
					}
				}

				ft := grm.WriterAtLine(ral)
				ts := bytes.TrimSuffix(ft.Bytes(), []byte("\n"))
				dists0 = append(dists0, ts)
				d.Reset()
			}
			dists0 = append(dists0, v)
			continue
		}

		d.Write(v)
		d.WriteByte('\n')
	}

	return bytes.Join(dists0, n)
}
