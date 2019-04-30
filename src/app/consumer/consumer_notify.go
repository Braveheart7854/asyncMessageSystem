package main

import (
	"app/common"
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/streadway/amqp"
	"log"
	"strconv"
)

var db *sql.DB

func init()  {
	var err error
	db, err = sql.Open("mysql", "root:hello123@tcp(127.0.0.1:3306)/test")
	common.FailOnError(err,"")
	db.SetMaxOpenConns(2000)
	db.SetMaxIdleConns(1000)
	//db.SetConnMaxLifetime(3)
	db.Ping()
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
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

	//if len(os.Args) < 2 {
	//	log.Printf("Usage: %s [binding_key]...", os.Args[0])
	//	os.Exit(0)
	//}
	//for _, s := range os.Args[1:] {
	//	log.Printf("Binding queue %s to exchange %s with routing key %s", q.Name, "logs_topic", s)
		err = ch.QueueBind(
			q.Name,       // queue name
			common.RouteKeyNotice,            // routing key
			common.ExchangeNameNotice, // exchange
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

	go func() {
		for d := range msgs {

			orderSn := common.MD5(string(d.Body)+"notify")
			//log.Printf(" [x] %s", orderSn)
			//log.Printf(" [x] %s", d.Body)

			var notice map[string] interface{}
			err := json.Unmarshal([]byte(d.Body),&notice)

			if err != nil{
				log.Println(err)
				if common.LogErrorJobs(db,orderSn,string(d.Body),"notify"){
					d.Ack(false)
				}else{
					d.Nack(false,true)
				}
				continue
			}

			index := common.GetHaseValue(int(notice["uid"].(float64)))
			table := "notice_" + strconv.Itoa(index)

			smtp,err := db.Prepare("insert into "+table+" (uid,type,data,create_time) values (?,?,?,?)")
			if err != nil {
				log.Println(err)
				if common.LogErrorJobs(db,orderSn,string(d.Body),"notify"){
					d.Ack(false)
				}else{
					d.Nack(false,true)
				}
				continue
			}
			result,err := smtp.Exec(notice["uid"],notice["type"],notice["data"],notice["createTime"])
			if err != nil {
				log.Println(err)
				if common.LogErrorJobs(db,orderSn,string(d.Body),"notify"){
					d.Ack(false)
				}else{
					d.Nack(false,true)
				}
				continue
			}
			lastId,err := result.LastInsertId()
			if common.IsEmpty(lastId) || err != nil {

				log.Println(err)
				if common.LogErrorJobs(db,orderSn,string(d.Body),"notify"){
					d.Ack(false)
				}else{
					d.Nack(false,true)
				}
				continue
			}else{
				smtpuser,_ := db.Prepare("update users set notification_count=notification_count+1 where id=?")
				res,_ := smtpuser.Exec(notice["uid"])
				_,_ = res.RowsAffected()
				smtpuser.Close()
			}
			d.Ack(false)
			smtp.Close()
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
