package router

import (
	"asyncMessageSystem/app/controller/producer"
	"github.com/kataras/iris"
)

func UrlPath(app *iris.Application){
	message := app.Party("/api/msg")
	handler := producer.Produce{}
	message.Post("/product", handler.Notify)
	message.Post("/read", handler.Read)


	user := app.Party("/api/user")
	user.Post("/product", handler.Notify)
}