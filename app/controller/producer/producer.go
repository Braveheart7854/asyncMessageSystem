package producer

import (
	"asyncMessageSystem/app/common"
	"asyncMessageSystem/app/model"
	"github.com/Braveheart7854/rabbitmqPool"
	"github.com/kataras/iris"
	"strconv"
	"time"
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
	Uid int64 `json:"uid"`
	Type int64 `json:"type"`
	Data string `json:"data"`
	CreateTime string `json:"createTime"`
}

type ReturnJson struct {
	Code int `json:"code"`
	Msg string `json:"msg"`
	Data map[string]interface{} `json:"data"`
}

/**
 新增消息
 */
func (P *Produce) Notify(ctx iris.Context)  {

	uid    := ctx.PostValueInt64Default("uid",0)
	n_type := ctx.PostValueInt64Default("type",0)
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

/**
 消息标记为已读
 */
func (P *Produce) Read(ctx iris.Context) {
	uid    := ctx.PostValueInt64Default("uid",0)
	n_type := ctx.PostValueInt64Default("type",common.TYPE_LIKE)
	data   := ctx.PostValueDefault("data","")
	createTime   := ctx.PostValueDefault("time",time.Now().Format("2006-01-02 15:04:05"))

	var noticeData Notice
	noticeData.Uid = uid
	noticeData.Type = n_type
	noticeData.Data = data
	noticeData.CreateTime = createTime

	go rabbitmqPool.AmqpServer.PutIntoQueue(common.ExchangeNameRead,common.RouteKeyRead,noticeData)

	//log.Printf("%d %d %s",uid,n_type,data)
	_,_ = ctx.JSON(ReturnJson{Code:10000,Msg:"success",Data: map[string]interface{}{"uid":uid,"type":n_type,"data":data}})
	return
}

/**
 消息列表
 */
func (P *Produce) List(ctx iris.Context){
	uid,_ := strconv.Atoi(ctx.FormValueDefault("uid","0"))
	typ,_ := strconv.Atoi(ctx.FormValueDefault("type","0"))
	page,_ := strconv.Atoi(ctx.FormValueDefault("page","1"))

	NoticeModel := new(model.Notice)
	list := NoticeModel.GetListByUid(uid,typ,page)
	unread := NoticeModel.CountUnReadByUid(uid,typ)
	_,_ = ctx.JSON(ReturnJson{Code:10000,Msg:"success",Data: map[string]interface{}{"list":list,"unread":unread}})
	return
}