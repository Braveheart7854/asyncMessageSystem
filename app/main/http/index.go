package main

import (
	"asyncMessageSystem/app/config"
	"asyncMessageSystem/app/middleware"
	"asyncMessageSystem/app/router"
	"github.com/kataras/iris"
	"net/http"
	"os"
	"time"
)

func main() {
	app := iris.New()
	//app.Logger().SetLevel("debug")

	middleware.LoadRabbitmq()

	router.Handler(app)

	srv := &http.Server{
		Addr:config.ServerAddr + ":3333",
		ReadTimeout: 3 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	if err := app.Run(iris.Server(srv)); err != nil{
		println(err.Error())
		os.Exit(0)
	}
}