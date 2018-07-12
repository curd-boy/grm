package format

import (
	"testing"
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
	t.Log(string(dd))
}
