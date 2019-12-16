package router

import (
	"asyncMessageSystem/app/controller/producer"
	user2 "asyncMessageSystem/app/controller/user"
	"github.com/kataras/iris"
)

func UrlPath(app *iris.Application){
	message := app.Party("/api/msg")
	handler := producer.Produce{}
	message.Post("/product", handler.Notify)
	message.Post("/read", handler.Read)
	message.Get("/list", handler.List)


	user := app.Party("/api/user")
	userHandler := user2.User{}
	user.Get("/info", userHandler.UserInfo)
}