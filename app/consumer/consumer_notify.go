package main

import (
	"asyncMessageSystem/app/common"
	"asyncMessageSystem/app/config"
	"asyncMessageSystem/app/controller/producer"
	"asyncMessageSystem/app/middleware"
	"asyncMessageSystem/app/model"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/streadway/amqp"
	"log"
)

func init()  {
	middleware.InitMysql()
}

func main() {
	conn, err := amqp.Dial(config.AmqpUrl)
	common.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	common.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		common.ExchangeNameNotice, // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	common.FailOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		common.QueueNameNotice,    // name
		true, // durable
		false, // delete when unusedt.msg.ext.msg.ext.msg.ext.msg.ex
		false,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	common.FailOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,       // queue name
		common.RouteKeyNotice,            // routing key
		common.ExchangeNameNotice, // exchange
		false,
		nil)
	common.FailOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	common.FailOnError(err, "Failed to register a consumer")

	forever := make(chan bool)
	failedQueues := new(model.FailedQueues)
	noticeModel := new(model.Notice)
	userModel := new(model.User)

	go func() {
		for d := range msgs {

			orderSn := common.MD5(string(d.Body)+"notify")

			var notice = new(producer.Notice)
			err := json.Unmarshal([]byte(string(d.Body)),notice)

			//fmt.Println(notice)

			if err != nil{
				log.Println("json.Unmarshal error: ",err.Error())
				if failedQueues.LogErrorJobs(orderSn,string(d.Body),"notify"){
					d.Ack(false)
				}else{
					d.Nack(false,true)
				}
				continue
			}

			table := noticeModel.TableName(notice.Uid)
			exist,errIsEx := noticeModel.IsExistNotice(table,orderSn)

			if errIsEx != nil{
				log.Println("noticeModel.IsExistNotice error: ",errIsEx.Error())
				if failedQueues.LogErrorJobs(orderSn,string(d.Body),"notify"){
					d.Ack(false)
				}else{
					d.Nack(false,true)
				}
				continue
			}
			if exist {
				d.Ack(false)
				continue
			}

			add, adderr := noticeModel.AddNotice(table,orderSn,notice.Uid,notice.Type,notice.Data,notice.CreateTime)
			if !add || adderr != nil {
				log.Println("noticeModel.AddNotice error: ",adderr.Error())

				if failedQueues.LogErrorJobs(orderSn,string(d.Body),"notify"){
					d.Ack(false)
				}else{
					d.Nack(false,true)
				}
				continue
			}else{
				_,_ = userModel.IncryNotifyCount(notice.Uid)
			}
			d.Ack(false)

		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
