package producer

import (
	"fmt"
	"github.com/Braveheart7854/rabbitmqPool"
	"github.com/kataras/iris"
	"log"
	"time"
	"wxforum_server/app/common"
)

type Produce struct {}

type Producer interface {
	Notify(ctx iris.Context)
	Read(ctx iris.Context)

}

//const (
//	EXCHANGE_NOTICE = "exchange_wxforum_notice"
//	ROUTE_NOTICE    = "route_wxforum_notice"
//)

type Notice struct {
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
	//common.Log("./log3.txt","test")

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

	var noticeData Notice
	noticeData.Uid = uid
	noticeData.Type = n_type
	noticeData.Data = data
	noticeData.CreateTime = createTime

	//common.Log("./log1.txt",data)
	go rabbitmqPool.AmqpServer.PutIntoQueue(common.ExchangeNameNotice,common.RouteKeyNotice,noticeData)
	//common.Log("./log2.txt",data)

	//log.Printf("%d %d %s",uid,n_type,data)
	ctx.JSON(ReturnJson{Code:10000,Msg:"success",Data: map[string]interface{}{"uid":uid,"type":n_type,"data":data}})
	return
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

	var noticeData Notice
	noticeData.Uid = uid
	noticeData.Type = n_type
	noticeData.Data = data
	noticeData.CreateTime = createTime

	go rabbitmqPool.AmqpServer.PutIntoQueue(common.ExchangeNameRead,common.RouteKeyRead,noticeData)

	//log.Printf("%d %d %s",uid,n_type,data)
	ctx.JSON(ReturnJson{Code:10000,Msg:"success",Data: map[string]interface{}{"uid":uid,"type":n_type,"data":data}})
	return
}

type quizzes struct {
	Qeid       int    `json:"qeid"`
	CreateTime string `json:"createTime"`
}
func (P * Produce) QuizzesPrize(ctx iris.Context){
	qeid    := ctx.PostValueIntDefault("qeid",0)
	createTime   := ctx.PostValueDefault("time",time.Now().Format(common.LAYOUT_STYLE))

	var data quizzes
	data.Qeid = qeid
	data.CreateTime = createTime

	//QueueService := new(Service)
	_,err := rabbitmqPool.AmqpServer.PutIntoQueue(common.ExchangeNameQuizzesEvent,common.RouteKeyQuizzesEvent,data)
	if err != nil {
		fmt.Println(err)
		common.Logger.Notice(data)
		ctx.JSON(ReturnJson{Code:common.FAILED,Msg:"failed",Data: map[string]interface{}{"qeid":qeid,"data":data}})
		return
	}
	//log.Printf("%d %d %s",uid,n_type,data)
	ctx.JSON(ReturnJson{Code:common.SUCCESS,Msg:"success",Data: map[string]interface{}{"qeid":qeid,"data":data}})
	return
}