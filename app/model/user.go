package model

import (
	"asyncMessageSystem/app/common"
	"asyncMessageSystem/app/model/db"
	"log"
	"time"
	"xorm.io/core"
)

type User struct {
	Id              uint64    `xorm:"unsigned uint64(11) pk autoincr not null id" json:"id"`
	Account         string    `xorm:"varchar(32) not null default '' index account" json:"account"`
	Password        string    `xorm:"varchar(32) not null default '' password" json:"-"`
	Nick            string    `xorm:"varchar(32) not null default '' nick" json:"nick"`
	NotificationCount uint64  `xorm:"uint64(20) not null default 0 notification_count" json:"notification_count"`
	Status          int8      `xorm:"tinyint(1) not null default 0 status" json:"status"`
	CreatedAt       time.Time `xorm:"datetime not null created created_at" json:"created_at"`
	UpdatedAt       time.Time `xorm:"datetime not null updated updated_at" json:"-"`
}

func (u *User) TableName()string{
	return "user"
}

type UserExecPrepare struct {
	SelectPrepare *core.Stmt
	InsertPrepare *core.Stmt
	UpdatePrepare *core.Stmt
}

var UserPrepare = new(UserExecPrepare)

func (u *User) InitPrepare(){
	var err error
	//UserPrepare.SelectPrepare,err = db.DB.DB().Prepare("select id from ? where order_sn = ?")
	//if err != nil{
	//	log.Panic(err)
	//}
	//UserPrepare.InsertPrepare,err = db.DB.DB().Prepare("insert into ? (order_sn,uid,type,data,create_time) values (?,?,?,?,?)")
	//if err != nil{
	//	log.Panic(err)
	//}
	UserPrepare.UpdatePrepare,err = db.DB.DB().Prepare("update user set notification_count=notification_count+1 and updated_at=? where id=?")
	if err != nil{
		log.Panic(err)
	}
}

func (u *User) GetUserInfoByUid(uid uint64)map[string]interface{}{
	var user = User{Id:uid}
	has,err := db.DB.Get(&user)
	if err != nil{
		return nil
	}
	if has == false{
		return nil
	}
	return map[string]interface{}{"nick":user.Nick,"notification_count":user.NotificationCount}
}

func (u *User) UpdateUserByUid(uid uint64)(bool,error){
	timeStr := time.Now().Format(common.LAYOUT_STYLE)
	result,err := UserPrepare.UpdatePrepare.Exec(timeStr,uid)
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