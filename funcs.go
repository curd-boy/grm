package grm

import (
	"errors"
	"fmt"
	"strings"
)

var Funcs = map[string]interface{}{
	"Commas": Commas,
	"Equal":  Equal,
	"Like":   Like,
	"And":    And,
	"Or":     Or,
	"Where":  Where,
	"Having": Having,
}

func Commas(a []string) string {
	return strings.Join(a, ", ")
}

func Equal(k string, v ...interface{}) (string, error) {
	ss := make([]string, 0, len(v))
	for _, v0 := range v {
		ss = append(ss, escapeParamsString(v0))
	}
	switch len(ss) {
	case 0:
		return "", errors.New("参数太少")
	case 1:
		return fmt.Sprintf(" %s = %s ", k, ss[0]), nil
	default:
		return fmt.Sprintf(" %s IN (%s) ", k, strings.Join(ss, ", ")), nil
	}
}

func Like(k string, v interface{}) (string, error) {
	return fmt.Sprintf(" %s LIKE %s ", k, escapeParamsString(v)), nil
}

func And(v []interface{}) string {
	ss := make([]string, 0, len(v))
	for _, v0 := range v {
		ss = append(ss, fmt.Sprintf("(%v)", v0))
	}
	return strings.Join(ss, " AND\n ")
}

func Or(v []interface{}) string {
	ss := make([]string, 0, len(v))
	for _, v0 := range v {
		ss = append(ss, fmt.Sprintf("(%v)", v0))
	}
	return strings.Join(ss, " OR\n ")
}

func Where(s string) string {
	if len(s) == 0 {
		return ""
	}
	return fmt.Sprintf("WHERE %s", s)
}

func Having(s string) string {
	if len(s) == 0 {
		return ""
	}
	return fmt.Sprintf("HAVING %s", s)
}
