package mysql

import (
	"asyncMessageSystem/app/common"
	."asyncMessageSystem/app/config"
	"asyncMessageSystem/app/model"
	"asyncMessageSystem/app/model/db"
	"github.com/go-xorm/xorm"
	"log"
	"time"
)

//func init(){
func Init(){
	InitDB()
	InitMigrate()
	InitPrepare()
}

func InitDB(){
	engine, err := xorm.NewEngine("mysql", Conf.Mysql.Dsn)
	common.FailOnError(err,"数据库连接池创建失败！")
	engine.SetMaxOpenConns(Conf.Xorm.MaxOpenConns)
	engine.SetMaxIdleConns(Conf.Xorm.MaxIdleConns)
	engine.SetConnMaxLifetime(Conf.Xorm.ConnMaxLifetime)

	engine.TZLocation, _ = time.LoadLocation(Conf.Xorm.TimeZone)
	engine.ShowSQL(Conf.Xorm.Debug) //调试用
	db.DB = engine
}

func InitMigrate(){
	notice := model.Notice{}
	for i:=1;i<= model.TableNum ;i++  {
		tablename := notice.TableName(uint64(i))
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

	createTable(model.User{},new(model.User).TableName())
	createTable(model.FailedQueues{},new(model.FailedQueues).TableName())
}

func createTable(table interface{},tableName string){
	res,err := db.DB.IsTableExist(table)
	if err !=nil{
		log.Panic(err.Error())
	}
	if res == false {
		err = db.DB.Charset("utf8mb4").CreateTable(table)
		if err != nil {
			log.Panic(err.Error())
		}
		err = db.DB.CreateIndexes(table)
		if err != nil {
			log.Panic(err.Error())
		}
		println("Created table ",tableName)
	}
}

func InitPrepare(){
	new(model.FailedQueues).InitPrepare()
	new(model.Notice).InitPrepare()
	new(model.User).InitPrepare()
}