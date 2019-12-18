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
		common.ExchangeNameRead, // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	common.FailOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		common.QueueNameRead,    // name
		true, // durable
		false, // delete when unusedt.msg.ext.msg.ext.msg.ext.msg.ex
		false,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	common.FailOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,       // queue name
		common.RouteKeyRead,            // routing key
		common.ExchangeNameRead, // exchange
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

			orderSn := common.MD5(string(d.Body) + "read")

			var read = new(producer.Notice)
			err := json.Unmarshal([]byte(string(d.Body)),read)

			if err != nil{
				log.Println(err)
				if failedQueues.LogErrorJobs(orderSn,string(d.Body),"read"){
					d.Ack(false)
				}else{
					d.Nack(false,true)
				}
				continue
			}

			table := new(model.Notice).TableName(read.Uid)

			rowsCount,errUpdate := noticeModel.UpdateNotice(table,read.Uid,read.Type)
			if errUpdate != nil{
				log.Println("noticeModel.UpdateNotice error: ",errUpdate.Error())
				if failedQueues.LogErrorJobs(orderSn,string(d.Body),"read"){
					d.Ack(false)
				}else{
					d.Nack(false,true)
				}
				continue
			}
			if rowsCount > 0 {
				decryCount,decryErr := userModel.DecryNotifyCount(read.Uid,rowsCount)
				if decryErr != nil {
					log.Println("userModel.DecryNotifyCount error: ",decryErr.Error())
				}
				if decryCount == 0{
					_,clearErr := userModel.ClearNotifyCount(read.Uid)
					if clearErr != nil {
						log.Println("userModel.ClearNotifyCount error: ",clearErr.Error())
					}
				}
			}
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
