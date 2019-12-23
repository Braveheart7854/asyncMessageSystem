package main

import (
	"asyncMessageSystem/app/config"
	_ "asyncMessageSystem/app/middleware/mysql"
	_ "asyncMessageSystem/app/middleware/rabbitmq"
	//_ "asyncMessageSystem/app/middleware/redis"
	"asyncMessageSystem/app/router"
	"github.com/kataras/iris"
	"net/http"
	"os"
)

func main() {

	app := iris.New()
	//app.Logger().SetLevel("debug")

	router.Handler(app)

	srv := &http.Server{
		Addr:config.ServerAddr,
		ReadTimeout: config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
	}
	if err := app.Run(iris.Server(srv)); err != nil{
		println(err.Error())
		os.Exit(0)
	}
}