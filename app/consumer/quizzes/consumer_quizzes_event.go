package main

import (
	"github.com/Braveheart7854/rabbitmqPool"
	"wxforum_server/app/common"
	"wxforum_server/app/config"
	"wxforum_server/app/producer"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/streadway/amqp"
	"log"
	"math"
	"strconv"
	"time"
)

var db *sql.DB

func init()  {
	var err error
	db, err = sql.Open("mysql", config.MysqlDataSource)
	common.FailOnError(err,"")
	db.SetMaxOpenConns(2000)
	db.SetMaxIdleConns(1000)
	db.SetConnMaxLifetime(9)
	db.Ping()

	rabbitmqPool.AmqpServer = rabbitmqPool.Service{
		AmqpUrl:config.AmqpUrl,
	}
	rabbitmqPool.InitAmqp()
}

type data struct {
	Type string `json:"type"`
	Content string `json:"content"`
	Reason string `json:"reason"`
}

func main() {
	defer db.Close()

	conn, err := amqp.Dial(config.AmqpUrl)
	common.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	common.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		common.ExchangeNameQuizzesEvent, // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	common.FailOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		common.QueueNameQuizzesEvent,    // name
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
		common.RouteKeyQuizzesEvent,            // routing key
		common.ExchangeNameQuizzesEvent, // exchange
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

	var noticeData producer.Notice
	var data data
	go func() {
		for d := range msgs {

			orderSn := common.MD5(string(d.Body)+"quizzes_event")
			//log.Printf(" [x] %s", orderSn)
			//log.Printf(" [x] %s", d.Body)

			var notice map[string] interface{}
			err := json.Unmarshal([]byte(d.Body),&notice)

			if err != nil{
				log.Println(err)
				if common.LogErrorJobs(db,orderSn,string(d.Body),"quizzes_event"){
					d.Ack(false)
				}else{
					d.Nack(false,true)
				}
				continue
			}

			var (
				id   int
				player1 int
				player2 int
				player1_coin int
				player2_coin int
				winer  int
				title  string
				)
			_ = db.QueryRow("select id,player1,player2,player1_coin,player2_coin,winer,title from quizzes_events where id=? and status=3",notice["qeid"]).
				Scan(&id,&player1,&player2,&player1_coin,&player2_coin,&winer,&title)
			if common.IsEmpty(id) {
				d.Ack(false)
				continue
			}
			totalCoin := player1_coin + player2_coin

			var userCount int
			_ = db.QueryRow("select count(id) from quizzes_records where qe_id = ? and status = 1", notice["qeid"]).Scan(&userCount)
			if userCount == 0 {
				d.Ack(false)
				continue
			}

			limit := 500
			for i := 1;i <= int(math.Ceil(float64(userCount)/float64(limit))) ;i++  {
				go func(i int) {
					var qrid int
					var uid int
					var player int
					var bet_coin int
					rows,_ := db.Query("select id,uid,player,bet_coin from quizzes_records where qe_id = ? and status = 1 limit ? offset ?",notice["qeid"],limit,(i-1)*limit)
					for rows.Next() {
						err := rows.Scan(&qrid,&uid, &player,&bet_coin)
						if err != nil {
							fmt.Println(err)
							continue
						}

						//打平，退回投注
						if winer == 0 {
							tx1,err := db.Begin()
							if err != nil {
								fmt.Println(err)
								continue
							}
							res,err := db.Exec("update quizzes_records set status=3 where id=?",qrid)
							if err != nil {
								fmt.Println(err)
								tx1.Rollback()
								continue
							}
							countAffect,err := res.RowsAffected()
							if countAffect ==0 || err!=nil {
								fmt.Println(err)
								tx1.Rollback()
								continue
							}

							resUser,err := db.Exec("update users set coin=coin+? where id=?",bet_coin,uid)
							if err != nil {
								fmt.Println(err)
								tx1.Rollback()
								continue
							}
							uCountAffect,err := resUser.RowsAffected()
							if uCountAffect ==0 || err!=nil {
								fmt.Println(err)
								tx1.Rollback()
								continue
							}
							coinLog,err := db.Exec("insert into coin_logs (user_id,type,number,content) values (?,?,?,?)",uid,1,bet_coin,title+" 平局，退回投注")
							if err != nil {
								fmt.Println(err)
								tx1.Rollback()
								continue
							}
							coinlogId,err := coinLog.LastInsertId()
							if coinlogId ==0 || err!=nil {
								fmt.Println(err)
								tx1.Rollback()
								continue
							}
							tx1.Commit()

							data.Type = "adminMsg"
							data.Content = "你获得"+ strconv.Itoa(bet_coin) +"咸豆奖励"
							data.Reason = title+" 平局，退回投注"
							jdata,_ := json.Marshal(data)

							noticeData.Uid = uid
							noticeData.Type = common.TYPE_SYSTEM
							noticeData.Data = string(jdata)
							noticeData.CreateTime = time.Now().Format(common.LAYOUT_STYLE)
							rabbitmqPool.AmqpServer.PutIntoQueue(common.ExchangeNameNotice,common.RouteKeyNotice,noticeData)
							continue
						}

						if player == winer {
							var prizeCoin int
							if player == player1 {
								if bet_coin > player1_coin {
									continue
								}
								prizeCoin = int(math.Ceil((float64(bet_coin)/float64(player1_coin))*float64(totalCoin)))
							}else if player == player2{
								if bet_coin > player2_coin {
									continue
								}
								prizeCoin = int(math.Ceil((float64(bet_coin)/float64(player2_coin))*float64(totalCoin)))
							}else {
								_,err := db.Exec("update quizzes_records set status=3 where id=?",qrid)
								if err != nil {
									fmt.Println(err)
								}
								continue
							}

							//给中奖者发奖
							tx,err := db.Begin()
							if err != nil {
								fmt.Println(err)
								continue
							}
							res,err := db.Exec("update quizzes_records set prize_coin=?, status=2 where id=?",prizeCoin,qrid)
							if err != nil {
								fmt.Println(err)
								tx.Rollback()
								continue
							}
							countAffect,err := res.RowsAffected()
							if countAffect ==0 || err!=nil {
								fmt.Println(err)
								tx.Rollback()
								continue
							}

							resUser,err := db.Exec("update users set coin=coin+? where id=?",prizeCoin,uid)
							if err != nil {
								fmt.Println(err)
								tx.Rollback()
								continue
							}
							uCountAffect,err := resUser.RowsAffected()
							if uCountAffect ==0 || err!=nil {
								fmt.Println(err)
								tx.Rollback()
								continue
							}

							coinLog,err := db.Exec("insert into coin_logs (user_id,type,number,content) values (?,?,?,?)",uid,1,prizeCoin,title+" 猜中，奖励咸豆")
							if err != nil {
								fmt.Println(err)
								tx.Rollback()
								continue
							}
							coinlogId,err := coinLog.LastInsertId()
							if coinlogId ==0 || err!=nil {
								fmt.Println(err)
								tx.Rollback()
								continue
							}
							tx.Commit()

							data.Type = "adminMsg"
							data.Content = "你获得"+ strconv.Itoa(prizeCoin) +"咸豆奖励"
							data.Reason = title+" 猜中，奖励咸豆"
							jdata,_ := json.Marshal(data)

							noticeData.Uid = uid
							noticeData.Type = common.TYPE_SYSTEM
							noticeData.Data = string(jdata)
							noticeData.CreateTime = time.Now().Format(common.LAYOUT_STYLE)
							rabbitmqPool.AmqpServer.PutIntoQueue(common.ExchangeNameNotice,common.RouteKeyNotice,noticeData)

						}else{

							_,err := db.Exec("update quizzes_records set status=3 where id=?",qrid)
							if err != nil {
								fmt.Println(err)
							}
						}

					}
					err = rows.Err()
					if err != nil {
						fmt.Println(err)
					}

					rows.Close()

				}(i)

			}

			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
