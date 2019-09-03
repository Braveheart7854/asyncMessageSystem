package main

import (
	"github.com/Braveheart7854/rabbitmqPool"
	"github.com/kataras/iris"
	"log"
	"asyncMessageSystem/app/config"
	"asyncMessageSystem/app/producer"
	"runtime/debug"
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

	app.Use(before)
	app.Post("/product", handler.Notify)
	app.Post("/read", handler.Read)

	app.Run(iris.Addr(config.ServerAddr + ":3333"))
}

func before(ctx iris.Context) {
	defer func() {
		msg := recover()
		if msg != nil {
			err := debug.Stack()
			log.Println(msg, "["+string(err)+"]")
			if !config.Debug {
				ctx.JSON(producer.ReturnJson{Code: 10001, Msg: "System is busy now!", Data: map[string]interface{}{}})
				return
			}else{
				strmsg := msg.(string)+"\r\n"
				bytemsg := []byte(strmsg)
				ctx.Write(append(bytemsg,err...))
				return
			}
		}
	}()
	ctx.Next() //继续执行下一个handler，在本例中是mainHandler。
}