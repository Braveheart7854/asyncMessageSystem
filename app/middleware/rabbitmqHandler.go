package middleware

import (
	"asyncMessageSystem/app/config"
	"github.com/Braveheart7854/rabbitmqPool"
)

func LoadRabbitmq(){
	rabbitmqPool.AmqpServer = rabbitmqPool.Service{
		AmqpUrl:config.AmqpUrl,
		ConnectionNum:10,
		ChannelNum:100,
	}
	rabbitmqPool.InitAmqp()
}
