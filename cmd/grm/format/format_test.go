package format

import (
	"testing"

	ffmt "gopkg.in/ffmt.v1"
)

var d = `{{define "GetStatus"}}
-- @Type Select
-- @Comm "get status"
-- @Req  KeySecretId       int     "KeySecretId"
-- @Req  ProxyId           int     "ProxyId"
-- @Resp Status            int     "Status"
--


{{end}}

`

func TestA(t *testing.T) {
	dd := Format([]byte(d))
	ffmt.Mark(string(dd))
}
