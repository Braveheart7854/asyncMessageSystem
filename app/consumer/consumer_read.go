package main

import (
	"asyncMessageSystem/app/common"
	"asyncMessageSystem/app/config"
	"asyncMessageSystem/app/middleware"
	"asyncMessageSystem/app/model"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/streadway/amqp"
	"log"
	"time"
)

var dbr *sql.DB

func init()  {
	var err error
	dbr, err = sql.Open("mysql", config.MysqlDataSource)
	common.FailOnError(err,"")
	dbr.SetMaxOpenConns(2000)
	dbr.SetMaxIdleConns(1000)
	dbr.SetConnMaxLifetime(9*time.Second)
	dbr.Ping()

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

	//if len(os.Args) < 2 {
	//	log.Printf("Usage: %s [binding_key]...", os.Args[0])
	//	os.Exit(0)
	//}
	//for _, s := range os.Args[1:] {
	//	log.Printf("Binding queue %s to exchange %s with routing key %s", q.Name, "logs_topic", s)
	err = ch.QueueBind(
		q.Name,       // queue name
		common.RouteKeyRead,            // routing key
		common.ExchangeNameRead, // exchange
		false,
		nil)
	common.FailOnError(err, "Failed to bind a queue")
	//}

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

	go func() {
		for d := range msgs {

			orderSn := common.MD5(string(d.Body) + "read")
			//log.Printf(" [x] %s", orderSn)
			//log.Printf(" [x] %s", d.Body)

			var read map[string] interface{}
			err := json.Unmarshal([]byte(d.Body),&read)

			if err != nil{
				log.Println(err)
				if failedQueues.LogErrorJobs(orderSn,string(d.Body),"read"){
					d.Ack(false)
				}else{
					d.Nack(false,true)
				}
				continue
			}

			//index := common.GetHaseValue(int(read["uid"].(float64)))
			//table := "notice_" + strconv.Itoa(index)
			table := new(model.Notice).TableName(uint64(read["uid"].(float64)))

			strType := common.NoticeType(int(read["type"].(float64)))

			rowsql := fmt.Sprintf("update %s set status=1 where uid=? and type in (%s) and status=0",table,strType)
			smtp,err := dbr.Prepare(rowsql)
			if err != nil {
				log.Println(err)
				if failedQueues.LogErrorJobs(orderSn,string(d.Body),"read"){
					d.Ack(false)
				}else{
					d.Nack(false,true)
				}
				continue
			}

			result,err := smtp.Exec(read["uid"])
			if err != nil {
				log.Println(err)
				if failedQueues.LogErrorJobs(orderSn,string(d.Body),"read"){
					d.Ack(false)
				}else{
					d.Nack(false,true)
				}
				continue
			}
			rowCount,err := result.RowsAffected()
			if common.IsEmpty(rowCount) && err != nil {

				log.Println(err)
				if failedQueues.LogErrorJobs(orderSn,string(d.Body),"read"){
					d.Ack(false)
				}else{
					d.Nack(false,true)
				}
				continue
			}else{
				smtpuserI,err := dbr.Prepare("update user set notification_count=notification_count-? where id=? and notification_count>=?")
				if err != nil {
					//log.Println(44444)
					log.Println(err)
				}
				res,err := smtpuserI.Exec(rowCount,read["uid"],rowCount)
				if err != nil {
					//log.Println(11111)
					log.Println(err)
				}
				affect,_ := res.RowsAffected()
				smtpuserI.Close()

				if common.IsEmpty(affect) {
					smtpuser,err := dbr.Prepare("update user set notification_count=0 where id=?")
					if err != nil {
						//log.Println(22222)
						log.Println(err)
					}
					_,err = smtpuser.Exec(read["uid"])
					if err != nil {
						//log.Println(33333)
						log.Println(err)
					}
					smtpuser.Close()
				}
			}
			d.Ack(false)
			smtp.Close()
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
