package redis

import (
	."asyncMessageSystem/app/config"
	"github.com/gomodule/redigo/redis"
	"log"
	"strconv"
)

var Cache = new(Instance)

//func init(){
func Init(){
	var err error
	password := redis.DialPassword(Conf.Redis.Password)
	database := redis.DialDatabase(Conf.Redis.Database)
	Cache.Conn,err = redis.Dial("tcp",Conf.Redis.Host+":"+ strconv.Itoa(Conf.Redis.Port) ,password,database)
	if err != nil {
		log.Panic(err.Error())
	}
}

type Instance struct {
	Conn redis.Conn
}

func (i *Instance) Set(key string,value interface{})(reply interface{},err error){
	reply,err = i.Conn.Do("set",key,value)
	return
}

func (i *Instance) SetEx(key string,expire int,value interface{})(reply interface{},err error){
	reply,err = i.Conn.Do("setEx",key,expire,value)
	return
}