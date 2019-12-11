package middleware

import (
	"asyncMessageSystem/app/config"
	"asyncMessageSystem/app/controller/producer"
	"github.com/kataras/iris"
	"log"
	"runtime/debug"
)

func PanicHandler(ctx iris.Context) {
	defer func() {
		msg := recover()
		if msg != nil {
			err := debug.Stack()
			log.Println(msg, "["+string(err)+"]")
			if !config.Debug {
				ctx.JSON(producer.ReturnJson{Code: 10001, Msg: "System is busy now!", Data: map[string]interface{}{}})
				return
			}else{
				strmsg := msg.(string)+"\r\n"
				bytemsg := []byte(strmsg)
				ctx.Write(append(bytemsg,err...))
				return
			}
		}
	}()
	ctx.Next() //继续执行下一个handler，在本例中是mainHandler。
}