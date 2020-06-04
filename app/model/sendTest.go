package model

import (
	"asyncMessageSystem/app/model/db"
	"log"
	"xorm.io/core"
)

type SendTest struct {
	Id              uint64    `xorm:"unsigned uint64(11) pk autoincr not null id" json:"id"`
	Time            string    `xorm:"varchar(32) not null default '' time" json:"time"`
	Order           int64     `xorm:"int64(11) not null default 0 order" json:"order"`
}

var InsertPre *core.Stmt

func (s SendTest) Prepare(){
	var err error
	InsertPre,err = db.DB.DB().Prepare("insert into send_test (`time`,`order`) values (?,?)")
	if err != nil{
		log.Panic(err)
	}
}

func (s SendTest) InsertSend(send SendTest) (bool, error) {
	//res,err := db.DB.Insert(send)
	result,err :=InsertPre.Exec(send.Time,send.Order)
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