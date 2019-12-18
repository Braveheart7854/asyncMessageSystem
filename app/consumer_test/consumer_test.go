package consumer_test

import (
	"asyncMessageSystem/app/middleware"
	"asyncMessageSystem/app/model"
	"testing"
)

func count(orderSn string, typ string)int{
	middleware.InitMysql()
	failedQueues := new(model.FailedQueues)
	count := failedQueues.CountFailedByOrderSn(orderSn,typ)
	return count
}

func TestCount(t *testing.T) {
	// 这里定义一个临时的结构体来存储测试case的参数以及期望的返回值
	for _, unit := range []struct {
		m        string
		n        string
		expected int
	}{
		{"29351b3fdf513806fc6a0ab1a721b7d3","notify", 1},
		{"f95e9bd13428f9209f139a9cf9cd800c","notify", 1},
	} {
		// 调用排列组合函数，与期望的结果比对，如果不一致输出错误
		if actually := count(unit.m, unit.n); actually != unit.expected {
			t.Errorf("combination: [%v], actually: [%v]", unit, actually)
		}
	}
}

func update(table string,uid uint64,typ int)(int64,error){
	middleware.InitMysql()
	noticeModel := new(model.Notice)
	res,err := noticeModel.UpdateNotice(table,uid,typ)
	return res,err
}

func TestUpdate(t *testing.T) {
	// 这里定义一个临时的结构体来存储测试case的参数以及期望的返回值
	for _, unit := range []struct {
		table      string
		uid        uint64
		typ        int
		expected   int64
		err        error
	}{
		{"th_notice_4",3,1, 1,nil},
	} {
		// 调用排列组合函数，与期望的结果比对，如果不一致输出错误
		if actually,err := update(unit.table, unit.uid,unit.typ); actually != unit.expected || err!=nil {
			t.Errorf("combination: [%v], actually: [%v,%v]", unit, actually,err)
		}
	}
}