package rabbitmq

import (
	"base/log"

	"github.com/streadway/amqp"
)

type MsgPubRouter struct {
	ch        *amqp.Channel
	conn      *amqp.Connection
	queueName string
}

func NewMsgPubRouter() *MsgPubRouter {
	return new(MsgPubRouter)
}

func (pub *MsgPubRouter) Start(addr string, queueName string) error {
	conn, err := amqp.Dial(addr)
	if err != nil {
		return err
	}
	//defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return err
	}
	//defer ch.Close()

	err = ch.ExchangeDeclare(
		queueName, // name
		"direct",  // type
		true,      // durable
		false,     // auto-deleted
		false,     // internal
		false,     // no-wait
		nil,       // arguments
	)
	log.PrintOnError(err, "Failed to declare an exchange")

	pub.ch = ch
	pub.conn = conn
	pub.queueName = queueName
	return nil
}

func (pub *MsgPubRouter) Pub(body []byte, table map[string]interface{}, routerKey string) error {
	err := pub.ch.Publish(
		pub.queueName, // exchange
		routerKey,     // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			Headers:      table,
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		})
	log.PrintOnError(err, "Failed to publish a message")
	return err
}

func (pub *MsgPubRouter) Close() {
	if pub.conn != nil {
		pub.conn.Close()
		pub.conn = nil
	}

	if pub.ch != nil {
		pub.ch.Close()
		pub.ch = nil
	}
}
