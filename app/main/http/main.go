package main

import (
	. "asyncMessageSystem/app/config"
	_ "asyncMessageSystem/app/middleware"
	"asyncMessageSystem/app/middleware/log"
	"asyncMessageSystem/app/router"
	"github.com/kataras/iris"
	"net/http"
	"os"
)

func main() {
	//go func() {
	//	for true {
	//		select {
	//		case <-time.After(500*time.Millisecond):
	//			fmt.Println("NumGoroutine : ",runtime.NumGoroutine())
	//		}
	//	}
	//}()

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