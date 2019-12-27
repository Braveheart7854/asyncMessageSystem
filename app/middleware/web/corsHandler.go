package web

import (
	"github.com/kataras/iris"
)

func CorsHandler(ctx iris.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")             //允许访问所有域
	ctx.Header("Access-Control-Allow-Headers", "Content-Type") //header的类型
	ctx.Header("content-type", "application/json")             //返回数据格式是json

	ctx.Next() //继续执行下一个handler，在本例中是mainHandler。
}