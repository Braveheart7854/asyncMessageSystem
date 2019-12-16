package middleware

import (
	"asyncMessageSystem/app/common"
	"asyncMessageSystem/app/config"
	"asyncMessageSystem/app/model"
	"asyncMessageSystem/app/model/db"
	"github.com/go-xorm/xorm"
	"log"
	"time"
)

func InitDB(){
	engine, err := xorm.NewEngine("mysql", config.MysqlDataSource)
	common.FailOnError(err,"数据库连接池创建失败！")
	engine.SetMaxOpenConns(config.MaxOpenConns)
	engine.SetMaxIdleConns(config.MaxIdleConns)
	engine.SetConnMaxLifetime(config.ConnMaxLifetime)

	engine.TZLocation, _ = time.LoadLocation(config.TimeZone)
	engine.ShowSQL(true) //调试用
	db.DB = engine
}

func InitMigrate(){
	notice := model.Notice{}
	for i:=1;i<=16 ;i++  {
		tablename := notice.TableName(i)
		res,err := db.DB.IsTableExist(tablename)
		if err !=nil{
			log.Panic(err.Error())
		}
		if res == false {
			err = db.DB.Charset("utf8mb4").Table(tablename).CreateTable(notice)
			if err != nil {
				log.Panic(err.Error())
			}
			err = db.DB.Table(tablename).CreateIndexes(notice)
			if err != nil {
				log.Panic(err.Error())
			}
			println("Created table ",tablename)
		}
	}

	user := model.User{}
	res,err := db.DB.IsTableExist(user)
	if err !=nil{
		log.Panic(err.Error())
	}
	if res == false {
		err = db.DB.Charset("utf8mb4").CreateTable(user)
		if err != nil {
			log.Panic(err.Error())
		}
		println("Created table ",user.TableName())
	}

	failedqueues := model.FailedQueues{}
	res,err = db.DB.IsTableExist(failedqueues)
	if err !=nil{
		log.Panic(err.Error())
	}
	if res == false {
		err = db.DB.Charset("utf8mb4").CreateTable(failedqueues)
		if err != nil {
			log.Panic(err.Error())
		}
		println("Created table ",failedqueues.TableName())
	}
}