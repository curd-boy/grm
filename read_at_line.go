package grm

import (
	"encoding/csv"
	"io"
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
