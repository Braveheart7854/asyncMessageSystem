package model

import "time"

type FailedQueues struct{
	Id              uint64    `xorm:"unsigned uint64(11) pk autoincr not null id" json:"id"`
	OrderSn         string    `xorm:"varchar(32) not null default '' index order_sn" json:"order_sn,omitempty"`
	Data            string    `xorm:"varchar(500) not null default '' data" json:"data"`
	Type            int       `xorm:"int(3) not null default 0 type" json:"type"`
	FailedCount     uint64    `xorm:"unsigned uint64(11) not null default 0 index failed_count" json:"failed_count"`
	CreatedAt       time.Time `xorm:"datetime not null created created_at" json:"created_at"`
	UpdatedAt       time.Time `xorm:"datetime not null updated updated_at" json:"updated_at"`
}

func (f *FailedQueues) TableName()string{
	return "failed_queues"
}