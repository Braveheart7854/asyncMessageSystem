package notice

import (
	"asyncMessageSystem/app/common"
	"asyncMessageSystem/app/config"
	producer2 "asyncMessageSystem/app/controller/producer"
	"context"
	"encoding/json"
	"github.com/Braveheart7854/rabbitmqPool"
	"log"
	"runtime/debug"
	"time"
)

type Producerpc struct {}

func before(reponse *ProducerResponse, errs error) {
	msg := recover()
	if msg != nil {
		err := debug.Stack()
		log.Println(msg, "["+string(err)+"]")
		if !config.Conf.Web.Debug {
			reponse = &ProducerResponse{Code: 10001,Msg:"System is busy now!",Data:[]byte{}}
			errs = nil
		} else {
			strmsg := msg.(string) + "\r\n"
			bytemsg := []byte(strmsg)
			reponse = &ProducerResponse{Code: 10001,Msg:"System is busy now!",Data:append(bytemsg, err...)}
			errs = nil
		}

	}
}

func (P *Producerpc) Notify(ctx context.Context, in *NoticeRequest) (reponse *ProducerResponse,errs error) {
	defer func() {
		msg := recover()
		if msg != nil {
			err := debug.Stack()
			log.Println(msg, "["+string(err)+"]")
			if !config.Conf.Web.Debug {
				reponse = &ProducerResponse{Code: common.FAILED,Msg:"System is busy now!",Data:[]byte{}}
			} else {
				strmsg := msg.(string) + "\r\n"
				bytemsg := []byte(strmsg)
				reponse = &ProducerResponse{Code: common.FAILED,Msg:"System is busy now!",Data:append(bytemsg, err...)}
			}
		}
	}()

	uid    := in.GetUid()
	n_type := in.GetType()
	data   := in.GetData()
	createTime   := in.GetCreateTime()
	if createTime == "" {
		createTime = time.Now().Format("2006-01-02 15:04:05")
	}

	var noticeData producer2.Notice
	noticeData.Uid = uint64(uid)
	noticeData.Type = int(n_type)
	noticeData.Data = data
	noticeData.CreateTime = createTime

	go rabbitmqPool.AmqpServer.PutIntoQueue(common.ExchangeNameNotice,common.RouteKeyNotice,noticeData)

	result := map[string]interface{}{"uid":uid,"type":n_type,"data":data}

	strResult,_ := json.Marshal(result)
	reponse = &ProducerResponse{Code: common.SUCCESS,Msg:"success",Data:strResult}
	return
}

func (P *Producerpc)Read(ctx context.Context, in *NoticeRequest) (reponse *ProducerResponse,errs error) {
	defer func() {
		msg := recover()
		if msg != nil {
			err := debug.Stack()
			log.Println(msg, "["+string(err)+"]")
			if !config.Conf.Web.Debug {
				reponse = &ProducerResponse{Code: common.FAILED,Msg:"System is busy now!",Data:[]byte{}}
			} else {
				strmsg := msg.(string) + "\r\n"
				bytemsg := []byte(strmsg)
				reponse = &ProducerResponse{Code: common.FAILED,Msg:"System is busy now!",Data:append(bytemsg, err...)}
			}
		}
	}()

	uid    := in.GetUid()
	n_type := in.GetType()
	data   := in.GetData()
	createTime   := in.GetCreateTime()
	if createTime == "" {
		createTime = time.Now().Format("2006-01-02 15:04:05")
	}

	var noticeData producer2.Notice
	noticeData.Uid = uint64(uid)
	noticeData.Type = int(n_type)
	noticeData.Data = data
	noticeData.CreateTime = createTime

	go rabbitmqPool.AmqpServer.PutIntoQueue(common.ExchangeNameRead,common.RouteKeyRead,noticeData)

	result := map[string]interface{}{"uid":uid,"type":n_type,"data":data}

	strResult,_ := json.Marshal(result)
	reponse = &ProducerResponse{Code: common.SUCCESS,Msg:"success",Data:strResult}
	return
}