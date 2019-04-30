package producer

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
)

type Service struct {
	Service ServiceInterface
}

type ServiceInterface interface {
	//ConnectMq()
	PutIntoQueue()
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
//
//func (P *Produce) GetMqInstance()*amqp.Channel {
//	P.MqInstance
//}
//
//func (P *Produce) ConnectMq()*amqp.Channel{
//	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
//	failOnError(err, "Failed to connect to RabbitMQ")
//	//defer conn.Close()
//
//	ch, err := conn.Channel()
//	failOnError(err, "Failed to open a channel")
//	//defer ch.Close()
//	return ch
//}

func (S *Service) PutIntoQueue(exchangeName string, routeKey string, notice notice) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exchangeName, // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	var data string

	body, err := json.Marshal(notice)
	if err != nil {
		log.Panic(err)
	}
	data = string(body)
	err = ch.Publish(
		exchangeName, // exchange
		routeKey,     //severityFrom(os.Args), // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(data),
		})

	failOnError(err, "Failed to publish a message")

	log.Printf(" [x] Sent %s", data)

}