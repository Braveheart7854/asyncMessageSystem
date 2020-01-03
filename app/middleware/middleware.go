package middleware

import (
	"asyncMessageSystem/app/middleware/log"
	"asyncMessageSystem/app/middleware/mysql"
	"asyncMessageSystem/app/middleware/rabbitmq"
	"asyncMessageSystem/app/middleware/redis"
	"runtime/debug"
)

var RequireMiddleware = map[string]func(){
	"mysql"    : mysql.Init,
	"rabbitmq" : rabbitmq.Init,
	"redis"    : redis.Init,
}

func init(){
	defer func() {
		msg := recover()
		if msg != nil {
			err := debug.Stack()
			log.MainLogger.Panic(msg.(string) + "["+string(err)+"]")
		}
	}()
	for _,value := range RequireMiddleware {
		value()
	}
}