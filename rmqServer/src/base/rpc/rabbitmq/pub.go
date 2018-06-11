package rabbitmq

import (
	"base/log"
	"github.com/streadway/amqp"
)

type MsgPub struct {
	ch        *amqp.Channel
	conn      *amqp.Connection
	queueName string
}

func NewMsgPub() *MsgPub {
	return new(MsgPub)
}

func (pub *MsgPub) Start(addr string, queueName string) error { //"amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(addr)
	if err != nil {
		return err
	}
	//defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	//defer ch.Close()

	q, err := ch.QueueDeclare(
		queueName, // name	"task_queue"
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return err
	}

	pub.ch = ch
	pub.queueName = q.Name
	return nil
}

func (pub *MsgPub) Pub(body []byte, table map[string]interface{}) error {
	err := pub.ch.Publish(
		"",            // exchange
		pub.queueName, // routing key
		false,         // mandatory
		false,
		amqp.Publishing{
			Headers:      table,
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         body,
		})
	log.Info("send: %v,%v", string(body), err)
	return err
}

func (pub *MsgPub) Close() {
	if pub.conn != nil {
		pub.conn.Close()
		pub.conn = nil
	}

	if pub.ch != nil {
		pub.ch.Close()
		pub.ch = nil
	}
}
