package main

import (
	."asyncMessageSystem/app/config"
	_ "asyncMessageSystem/app/middleware"
	"asyncMessageSystem/app/router"
	"github.com/kataras/iris"
	"log"
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
	}
	if err := app.Run(iris.Server(srv)); err != nil{
		log.Println(err.Error())
		os.Exit(0)
	}
}