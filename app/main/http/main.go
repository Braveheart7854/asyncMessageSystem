package main

import (
	"asyncMessageSystem/app/middleware/log"
	."asyncMessageSystem/app/config"
	_ "asyncMessageSystem/app/middleware"
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
		Addr:Conf.Web.ServerAddr,
		ReadTimeout: Conf.Web.ReadTimeout,
		WriteTimeout: Conf.Web.WriteTimeout,
		IdleTimeout: Conf.Web.IdleTimeout,
	}
	if err := app.Run(iris.Server(srv)); err != nil{
		log.MainLogger.Error(err.Error())
		os.Exit(0)
	}
}