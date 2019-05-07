package main

import (
	"app/config"
	"app/producer"
	"github.com/kataras/iris"
)

func main() {
	app := iris.New()

	handler := new(producer.Produce)
	app.Post("/product", handler.Notify)
	app.Post("/read", handler.Read)

	//channel := handler.ConnectMq()
	//app.Post("/product", func(context iris.Context) {
	//	handler.Notify(context,channel)
	//})

	app.Run(iris.Addr(config.ServerAddr + ":3333"))
}