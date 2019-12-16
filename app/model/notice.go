package model

import (
	"asyncMessageSystem/app/common"
	"asyncMessageSystem/app/model/db"
	"fmt"
	"strconv"
	"time"
)

type Notice struct {
	Id              uint64    `xorm:"unsigned uint64(11) pk autoincr not null id" json:"id"`
	OrderSn         string    `xorm:"varchar(32) not null default '' index order_sn" json:"order_sn,omitempty"`
	Uid             uint64    `xorm:"unsigned uint64(11) not null default 0 index uid" json:"uid,omitempty"`
	Type            int       `xorm:"int(3) not null default 0 type" json:"type"`
	Data            string    `xorm:"varchar(500) not null default '' data" json:"data"`
	Status          int8      `xorm:"tinyint(1) not null default 0 status" json:"status"`
	CreateTime      time.Time `xorm:"datetime not null create_time" json:"-"`
	CreatedAt       time.Time `xorm:"datetime not null created created_at" json:"-"`
	UpdatedAt       time.Time `xorm:"datetime not null updated updated_at" json:"updated_at"`
}

func (n *Notice) TableName(uid int) (table string){
	index := common.GetHaseValue(uid)
	table = "th_notice_" + strconv.Itoa(index)
	return
}

func (n *Notice) Insert(){
	i,err := db.DB.Table(n.TableName(111)).Insert(Notice{OrderSn:"ddd",CreateTime:time.Now()},Notice{OrderSn:"fff",CreateTime:time.Now()})
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("i: ",i)
}

func (n *Notice) Update(){
	i,err := db.DB.Table(n.TableName(111)).Update(Notice{OrderSn:"eee"},Notice{Id:1})
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("u: ",i)
}

func (n *Notice) GetListByUid(uid int,typ int,page int)(list []Notice){
	notices := &Notice{}
	rows,err := db.DB.Table(n.TableName(uid)).Select("id,type,data,status,updated_at").
		Where("uid = ? and type = ?",uid,typ).
		OrderBy("created_at desc").Limit(10,(page-1)*10).Rows(notices)
	defer rows.Close()
	if err != nil {
		return nil
	}
	for rows.Next() {
		_ = rows.Scan(notices)
		list = append(list,*notices)
	}
	return
}

func (n *Notice) CountUnReadByUid(uid int,typ int)(count int64){
	notice := Notice{}
	count,err := db.DB.Table(n.TableName(uid)).Where("uid = ? and type = ?",uid,typ).Count(&notice)
	if err != nil {
		return
	}
	return
}