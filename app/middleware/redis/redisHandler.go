package redis

import (
	"github.com/gomodule/redigo/redis"
	"log"
)

var Cache = new(Instance)

func init(){
	var err error
	password := redis.DialPassword("")
	database := redis.DialDatabase(1)
	Cache.Conn,err = redis.Dial("tcp","127.0.0.1:6379",password,database)
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