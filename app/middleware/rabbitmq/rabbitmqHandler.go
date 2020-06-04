package rabbitmq

import (
	. "asyncMessageSystem/app/config"
	//"asyncMessageSystem/app/middleware/rabbitmqPool"
	"github.com/Braveheart7854/rabbitmqPool"
)

//func init(){
func Init(){
	rabbitmqPool.AmqpServer = rabbitmqPool.Service{
		AmqpUrl:Conf.RabbitMq.Dsn,
		ConnectionNum:Conf.RabbitMq.ConnectionNum,
		ChannelNum:Conf.RabbitMq.ChannelNum,
	}
	rabbitmqPool.InitAmqp()
}
