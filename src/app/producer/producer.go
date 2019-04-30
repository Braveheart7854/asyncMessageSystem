package producer

import (
	"app/common"
	"github.com/kataras/iris"
	"log"
	"time"
)

type Produce struct {
	//MqInstance amqp.Channel
}

type Producer interface {
	Notify(ctx iris.Context)
	Read(ctx iris.Context)

}

//const (
//	EXCHANGE_NOTICE = "exchange_wxforum_notice"
//	ROUTE_NOTICE    = "route_wxforum_notice"
//)

type notice struct {
	Uid int `json:"uid"`
	Type int `json:"type"`
	Data string `json:"data"`
	CreateTime string `json:"createTime"`
}

type ReturnJson struct {
	Code int `json:"code"`
	Msg string `json:"msg"`
	Data map[string]interface{} `json:"data"`
}

func (P *Produce) Notify(ctx iris.Context)  {
	defer func() {
		msg := recover()
		if msg != nil {
			log.Printf("%s",msg)
			ctx.JSON(ReturnJson{Code:10001,Msg:"System is busy now!",Data: map[string]interface{}{}})
			return
		}
	}()

	uid    := ctx.PostValueIntDefault("uid",0)
	n_type := ctx.PostValueIntDefault("type",0)
	data   := ctx.PostValueDefault("data","")
	createTime   := ctx.PostValueDefault("time",time.Now().Format("2006-01-02 15:04:05"))

	var noticeData notice
	noticeData.Uid = uid
	noticeData.Type = n_type
	noticeData.Data = data
	noticeData.CreateTime = createTime

	QueueService := new(Service)
	QueueService.PutIntoQueue(common.ExchangeNameNotice,common.RouteKeyNotice,noticeData)

	//log.Printf("%d %d %s",uid,n_type,data)
	ctx.JSON(ReturnJson{Code:10000,Msg:"success",Data: map[string]interface{}{"uid":uid,"type":n_type,"data":data}})
}

func (P *Produce) Read(ctx iris.Context) {
	defer func() {
		msg := recover()
		if msg != nil {
			log.Printf("%s",msg)
			ctx.JSON(ReturnJson{Code:10001,Msg:"System is busy now!",Data: map[string]interface{}{}})
			return
		}
	}()

	uid    := ctx.PostValueIntDefault("uid",0)
	n_type := ctx.PostValueIntDefault("type",common.TYPE_LIKE)
	data   := ctx.PostValueDefault("data","")
	createTime   := ctx.PostValueDefault("time",time.Now().Format("2006-01-02 15:04:05"))

	var noticeData notice
	noticeData.Uid = uid
	noticeData.Type = n_type
	noticeData.Data = data
	noticeData.CreateTime = createTime

	QueueService := new(Service)
	QueueService.PutIntoQueue(common.ExchangeNameRead,common.RouteKeyRead,noticeData)

	//log.Printf("%d %d %s",uid,n_type,data)
	ctx.JSON(ReturnJson{Code:10000,Msg:"success",Data: map[string]interface{}{"uid":uid,"type":n_type,"data":data}})
}