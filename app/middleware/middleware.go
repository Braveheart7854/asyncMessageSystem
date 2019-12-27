package middleware

import (
	"asyncMessageSystem/app/middleware/mysql"
	"asyncMessageSystem/app/middleware/rabbitmq"
	"asyncMessageSystem/app/middleware/redis"
)

var RequireMiddleware = map[string]func(){
	"mysql"    : mysql.Init,
	"rabbitmq" : rabbitmq.Init,
	"redis"    : redis.Init,
}

func init(){
	for _,value := range RequireMiddleware {
		value()
	}
}