package rabbitmq

import (
	"asyncMessageSystem/app/config"
	"github.com/Braveheart7854/rabbitmqPool"
)

func init(){
	rabbitmqPool.AmqpServer = rabbitmqPool.Service{
		AmqpUrl:config.AmqpUrl,
		ConnectionNum:config.ConnectionNum,
		ChannelNum:config.ChannelNum,
	}
	rabbitmqPool.InitAmqp()
}
