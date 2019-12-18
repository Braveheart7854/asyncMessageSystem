package model

import (
	"asyncMessageSystem/app/common"
	"asyncMessageSystem/app/model/db"
	"log"
	"time"
	"xorm.io/core"
)

type FailedQueues struct{
	Id              uint64    `xorm:"unsigned uint64(11) pk autoincr not null id" json:"id"`
	OrderSn         string    `xorm:"varchar(32) not null default '' index order_sn" json:"order_sn,omitempty"`
	Data            string    `xorm:"varchar(500) not null default '' data" json:"data"`
	Type            string    `xorm:"varchar(50) not null default '' type" json:"type"`
	FailedCount     uint64    `xorm:"unsigned uint64(11) not null default 0 index failed_count" json:"failed_count"`
	CreatedAt       time.Time `xorm:"datetime not null created created_at" json:"created_at"`
	UpdatedAt       time.Time `xorm:"datetime not null updated updated_at" json:"updated_at"`
}

func (f *FailedQueues) TableName()string{
	return "failed_queues"
}

type ExecPrepare struct {
	SelectPrepare *core.Stmt
	InsertPrepare *core.Stmt
	UpdatePrepare *core.Stmt
}

var FailedQueuesPrepare = new(ExecPrepare)

func (f *FailedQueues) InitPrepare(){
	var err error
	FailedQueuesPrepare.SelectPrepare,err = db.DB.DB().Prepare("select failed_count from failed_queues where order_sn = ? and type = ?")
	if err != nil{
		log.Panic(err)
	}
	FailedQueuesPrepare.InsertPrepare,err = db.DB.DB().Prepare("insert into failed_queues (order_sn,data,failed_count,type,created_at,updated_at) values (?,?,?,?,?,?)")
	if err != nil{
		log.Panic(err)
	}
	FailedQueuesPrepare.UpdatePrepare,err = db.DB.DB().Prepare("update failed_queues set failed_count=failed_count+1 and updated_at=? where order_sn=? and type = ?")
	if err != nil{
		log.Panic(err)
	}
}

func (f *FailedQueues) CountFailedByOrderSn(orderSn string, typ string)(count int){
	_ = FailedQueuesPrepare.SelectPrepare.QueryRow(orderSn,typ).Scan(&count)
	return
}

func (f *FailedQueues) UpdateByOrderSn(orderSn string,typ string)(bool,error){
	result,err := FailedQueuesPrepare.UpdatePrepare.Exec(time.Now().Format(common.LAYOUT_STYLE),orderSn,typ)
	if err != nil {
		return false,err
	}
	res,e := result.RowsAffected()
	if e != nil {
		return false,e
	}
	if res > 0{
		return true,nil
	}
	return false,nil
}

func (f *FailedQueues) InsertOrder(orderSn string,data string,typ string)(bool,error){
	timeStr := time.Now().Format(common.LAYOUT_STYLE)
	result,err := FailedQueuesPrepare.InsertPrepare.Exec(orderSn,data,1,typ,timeStr,timeStr)
	if err != nil {
		return false,err
	}
	res,e := result.LastInsertId()
	if e != nil {
		return false,e
	}
	if res > 0{
		return true,nil
	}
	return false,nil
}

func (f *FailedQueues) LogErrorJobs(orderSn string,data string,typ string)(flag bool){
	failCount := f.CountFailedByOrderSn(orderSn,typ)
	if failCount >= 5 {
		flag = true
	}else if failCount >= 1 && failCount < 5 {
		_,err := f.UpdateByOrderSn(orderSn,typ)
		if err != nil {
			log.Println("LogErrorJobs UpdateByOrderSn error: ",err.Error())
		}
		flag = false
	}else {
		_,err := f.InsertOrder(orderSn,data,typ)
		if err != nil {
			log.Println("LogErrorJobs InsertOrder error: ",err.Error())
		}
		flag = false
	}
	return
}