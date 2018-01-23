package grm

import (
	"bytes"
	"encoding/csv"
	"io"
	"strings"

	ffmt "gopkg.in/ffmt.v1"
)

func ReadAtLine(r io.Reader) ([][]string, error) {
	read := csv.NewReader(r)
	read.Comma = ' '
	read.TrimLeadingSpace = true
	read.LazyQuotes = true
	read.FieldsPerRecord = -1
	d, err := read.ReadAll()
	if err != nil {
		return nil, err
	}
	return d, nil
}

func WriterAtLine(ff [][]string) *bytes.Buffer {

	for k1, v1 := range ff {
		for k2, v2 := range v1 {
			if strings.Contains(v2, " ") {
				ff[k1][k2] = `"` + v2 + `"`
			}
		}
	}

	buf := bytes.NewBuffer(nil)

	tt := ffmt.FmtTable(ff)
	for _, v := range tt {
		buf.WriteString(strings.TrimSpace(v))
		buf.WriteByte('\n')
	}
	return buf
}
