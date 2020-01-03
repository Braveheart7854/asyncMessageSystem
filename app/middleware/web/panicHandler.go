package web

import (
	"asyncMessageSystem/app/config"
	"asyncMessageSystem/app/controller/producer"
	log2 "asyncMessageSystem/app/middleware/log"
	"github.com/kataras/iris"
	"log"
	"runtime/debug"
)

func PanicHandler(ctx iris.Context) {
	defer func() {
		msg := recover()
		if msg != nil {
			err := debug.Stack()
			log2.MainLogger.Error(msg.(string) + "["+string(err)+"]")
			if config.Conf.Web.Debug {
				log.Println(msg, "["+string(err)+"]")
				strmsg := msg.(string)+"\r\n"
				bytemsg := []byte(strmsg)
				ctx.Write(append(bytemsg,err...))
				return
			}else{
				ctx.JSON(producer.ReturnJson{Code: 10001, Msg: "System is busy now!", Data: map[string]interface{}{}})
				return
			}
		}
	}()
	ctx.Next() //继续执行下一个handler，在本例中是mainHandler。
}