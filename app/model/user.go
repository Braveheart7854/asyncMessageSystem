package model

import (
	"asyncMessageSystem/app/model/db"
	"time"
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