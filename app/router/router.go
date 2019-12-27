package router

import (
	"asyncMessageSystem/app/config"
	"asyncMessageSystem/app/middleware/web"
	"github.com/kataras/iris"
)

func Handler(app *iris.Application){

	//加载中间件
	//奔溃恢复
	app.Use(web.PanicHandler)
	//跨域
	if config.Conf.Web.Debug == true{
		app.Use(web.CorsHandler)
	}

	//加载路由
	UrlPath(app)
}