package middleware

import (
	"asyncMessageSystem/app/middleware/log"
	"asyncMessageSystem/app/middleware/mysql"
	"asyncMessageSystem/app/middleware/rabbitmq"
	"asyncMessageSystem/app/middleware/redis"
)

var RequireMiddleware = map[string]func(){
	"logger"   : log.Init,
	"mysql"    : mysql.Init,
	"rabbitmq" : rabbitmq.Init,
	"redis"    : redis.Init,
}

func init(){
	for _,value := range RequireMiddleware {
		value()
	}
}