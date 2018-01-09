package grm

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

var (
	_Null     = []byte("NULL")
	_True     = []byte{'1'}
	_False    = []byte{'0'}
	_ZeroTime = []byte("'0000-00-00'")
)

func escapeParamsString(arg interface{}) string {
	return string(escapeParams(arg))
}

func escapeParams(arg interface{}) (buf []byte) {
	if arg == nil {
		buf = _Null
		return
	}

	switch v := arg.(type) {
	case uint64:
		buf = []byte(strconv.FormatUint(v, 10))
	case uint32:
		buf = []byte(strconv.FormatUint(uint64(v), 10))
	case uint16:
		buf = []byte(strconv.FormatUint(uint64(v), 10))
	case uint8:
		buf = []byte(strconv.FormatUint(uint64(v), 10))
	case uint:
		buf = []byte(strconv.FormatUint(uint64(v), 10))

	case int64:
		buf = []byte(strconv.FormatInt(v, 10))
	case int32:
		buf = []byte(strconv.FormatInt(int64(v), 10))
	case int16:
		buf = []byte(strconv.FormatInt(int64(v), 10))
	case int8:
		buf = []byte(strconv.FormatInt(int64(v), 10))
	case int:
		buf = []byte(strconv.FormatInt(int64(v), 10))

	case float64:
		buf = []byte(strconv.FormatFloat(v, 'g', -1, 64))
	case float32:
		buf = []byte(strconv.FormatFloat(float64(v), 'g', -1, 64))

	case bool:
		if v {
			buf = _True
		} else {
			buf = _False
		}
	case time.Time:
		if v.IsZero() {
			buf = _ZeroTime
		} else {
			buf = append(buf, '\'')
			buf = append(buf, v.Format(time.RFC3339Nano)...)
			buf = append(buf, '\'')
		}
	case []byte:
		if v == nil {
			buf = _Null
		} else {
			buf = append(buf, "_binary'"...)
			buf = append(buf, escapeBytesBackslash(v)...)
			buf = append(buf, '\'')
		}
	case string:
		buf = append(buf, '\'')
		buf = append(buf, escapeBytesBackslash([]byte(v))...)
		buf = append(buf, '\'')
	case fmt.Stringer:
		buf = append(buf, '\'')
		buf = append(buf, escapeBytesBackslash([]byte(v.String()))...)
		buf = append(buf, '\'')
	default:
		if val := reflect.ValueOf(arg); val.Kind() == reflect.Ptr && !val.IsNil() {
			return escapeParams(val.Elem().Interface())
		}
	}
	return buf
}

// escapeBytesBackslash escapes []byte with backslashes (\)
// This escapes the contents of a string (provided as []byte) by adding backslashes before special
// characters, and turning others into specific escape sequences, such as
// turning newlines into \n and null bytes into \0.
// https://github.com/mysql/mysql-server/blob/mysql-5.7.5/mysys/charset.c#L823-L932
func escapeBytesBackslash(v []byte) []byte {
	buf := make([]byte, 0, len(v)*2)

	for _, c := range v {
		switch c {
		case '\x00':
			buf = append(buf, '\\', '0')
		case '\n':
			buf = append(buf, '\\', 'n')
		case '\r':
			buf = append(buf, '\\', 'r')
		case '\x1a':
			buf = append(buf, '\\', 'Z')
		case '\'':
			buf = append(buf, '\\', '\'')
		case '"':
			buf = append(buf, '\\', '"')
		case '\\':
			buf = append(buf, '\\', '\\')
		default:
			buf = append(buf, c)
		}
	}

	return buf
}
