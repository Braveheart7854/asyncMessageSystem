package model

import (
	"asyncMessageSystem/app/common"
	"asyncMessageSystem/app/model/db"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"
	"xorm.io/core"
)

type Notice struct {
	Id              uint64    `xorm:"unsigned uint64(11) pk autoincr not null id" json:"id"`
	OrderSn         string    `xorm:"varchar(32) not null default '' index order_sn" json:"order_sn,omitempty"`
	Uid             uint64    `xorm:"unsigned uint64(11) not null default 0 index uid" json:"uid,omitempty"`
	Type            int       `xorm:"int(3) not null default 0 type" json:"type"`
	Data            string    `xorm:"varchar(500) not null default '' data" json:"data"`
	Status          int       `xorm:"tinyint(1) not null default 0 status" json:"status"`
	CreateTime      time.Time `xorm:"datetime not null create_time" json:"-"`
	CreatedAt       time.Time `xorm:"datetime not null created created_at" json:"-"`
	UpdatedAt       time.Time `xorm:"datetime not null updated updated_at" json:"updated_at"`
}

const TableNum = 16

func (n *Notice) GetHaseValue(uid uint64)uint64 {
	return (uid % TableNum)+1
}

func (n *Notice) TableName(uid uint64) (table string){
	index := n.GetHaseValue(uid)
	table = "th_notice_" + strconv.FormatUint(index,10)
	return
}

type NoticeExecPrepare struct {
	SelectPrepare *core.Stmt
	InsertPrepare *core.Stmt
	UpdatePrepare *core.Stmt
}

var NoticePrepare = make(map[string]NoticeExecPrepare,TableNum)

func (n *Notice) InitPrepare(){
	for i:=1;i<= TableNum ; i++ {
		tableName := n.TableName(uint64(i))
		selectPre,err := db.DB.DB().Prepare("select id from "+tableName+" where order_sn = ?")
		if err != nil{
			log.Panic(err)
		}
		insertPre,err := db.DB.DB().Prepare("insert into "+tableName+" (order_sn,uid,type,data,create_time,created_at,updated_at) values (?,?,?,?,?,?,?)")
		if err != nil{
			log.Panic(err)
		}
		updatePre,err := db.DB.DB().Prepare("update "+tableName+" set status=1 , updated_at=? where uid=? and type=? and status=0")
		if err != nil{
			log.Panic(err)
		}
		NoticePrepare[tableName] = NoticeExecPrepare{SelectPrepare:selectPre,InsertPrepare:insertPre,UpdatePrepare:updatePre}
	}
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

func (n *Notice) GetListByUid(uid uint64,typ int,page int)(list []Notice){
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

func (n *Notice) CountUnReadByUid(uid uint64,typ int)(count int64){
	notice := Notice{}
	count,_ = db.DB.Table(n.TableName(uid)).Where("uid = ? and type = ?",uid,typ).Count(&notice)
	return
}

func (n *Notice) IsExistNotice(table string,orderSn string)(flag bool,err error){
	var id int
	err = NoticePrepare[table].SelectPrepare.QueryRow(orderSn).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return false,nil
		}
		return false,err
	}
	if id >0 {
		return true,nil
	}else{
		return false,nil
	}
}

func (n *Notice) AddNotice(table string,orderSn string,uid uint64,typ int,data string,createTime string)(bool,error){
	timeStr := time.Now().Format(common.LAYOUT_STYLE)
	result,err := NoticePrepare[table].InsertPrepare.Exec(orderSn,uid,typ,data,createTime,timeStr,timeStr)
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

func (n *Notice) UpdateNotice(table string,uid uint64,typ int)(int64,error){
	result,err := NoticePrepare[table].UpdatePrepare.Exec(time.Now().Format(common.LAYOUT_STYLE),uid,typ)
	if err != nil {
		return 0,err
	}
	res,e := result.RowsAffected()
	if e != nil {
		return 0,e
	}
	return res,nil
}