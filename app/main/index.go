package main

import (
	"github.com/Braveheart7854/rabbitmqPool"
	"github.com/kataras/iris"
	"log"
	"wxforum_server/app/config"
	"wxforum_server/app/producer"
)

const(
	debug = true
)

func main() {
	app := iris.New()
	//app.Logger().SetLevel("debug")

	rabbitmqPool.AmqpServer = rabbitmqPool.Service{
		AmqpUrl:config.AmqpUrl,
		ConnectionNum:10,
		ChannelNum:100,
	}
	rabbitmqPool.InitAmqp()

	handler := new(producer.Produce)

	//conn, err := amqp.Dial(config.AmqpUrl)
	//common.FailOnError(err, "Failed to connect to RabbitMQ")
	//defer conn.Close()
	//handler.MqInstance = conn

	app.Post("/product", handler.Notify)
	app.Post("/read", handler.Read)
	app.Post("/quizzes/prize", before, handler.QuizzesPrize)

	//channel := handler.ConnectMq()
	//app.Post("/product", func(context iris.Context) {
	//	handler.Notify(context,channel)
	//})

	app.Run(iris.Addr(config.ServerAddr + ":3333"))
}

func before(ctx iris.Context) {
	defer func() {
		if !debug {
			msg := recover()
			if msg != nil {
				log.Printf("%s",msg)
				ctx.JSON(producer.ReturnJson{Code:10001,Msg:"System is busy now!",Data: map[string]interface{}{}})
				return
			}
		}
	}()
	ctx.Next() //继续执行下一个handler，在本例中是mainHandler。
}