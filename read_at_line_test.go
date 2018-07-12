package grm

import (
	"bytes"
	"fmt"
	"testing"
)

func TestReadAtLine(t *testing.T) {
	d := `-- @Type Select
-- @Req  LotteryCountRecordId int       "抽奖次数id"
-- @Req  GroupBy              []string  "分组"
-- @Req  OrderBy              []string  "排序"
-- @Req  Offset               int       "索引偏移"
-- @Req  Limit                int       "获取数量"
-- @Resp LotteryCountRecordId int       "抽奖次数id"
-- @Resp UserId               int       "用户id"
-- @Resp LotteryPoolId        int       "奖池id"           
-- @Resp LotteryCount         int       "已经抽奖的次数"
-- @Resp LotteryCountMax      int       "最大抽奖的次数"
-- @Resp CreateTime           time.Time "创建时间"
-- @Resp UpdateTime           time.Time "更新时间"`

	dd, err := ReadAtLine(bytes.NewBufferString(d))
	if err != nil {
		t.Fatal(err)
		return
	}
	fmt.Println(dd)
}
