package main

import (
	"asyncMessageSystem/app/config"
	"asyncMessageSystem/app/middleware"
	"asyncMessageSystem/app/router"
	"github.com/kataras/iris"
	"net/http"
	"os"
)

func main() {
	app := iris.New()
	//app.Logger().SetLevel("debug")

	middleware.InitDB()
	middleware.InitMigrate()

	middleware.InitRabbitmq()

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