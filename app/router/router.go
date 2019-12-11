package router

import (
	"asyncMessageSystem/app/middleware"
	"github.com/kataras/iris"
)

func Handler(app *iris.Application){

	//加载中间件
	app.Use(middleware.PanicHandler)

	//加载路由
	UrlPath(app)
}