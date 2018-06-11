package rabbitmq

import (
	"base/log"
	"base/rpc"
	"github.com/streadway/amqp"
)

type MsgSub struct {
	ch   *amqp.Channel
	conn *amqp.Connection

	readChan <-chan amqp.Delivery
	stopSign chan byte
}

func NewMsgSub() *MsgSub {
	return new(MsgSub)
}

func (sub *MsgSub) Start(addr string, queueName string) error {
	var ok bool
	//异常退出释放
	defer func() {
		if !ok {
			sub.onClose()
		}
	}()

	conn, err := amqp.Dial(addr)
	if err != nil {
		return err
	}
	sub.conn = conn
	//defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	//defer ch.Close()
	sub.ch = ch

	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return err
	}

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return err
	}

	sub.readChan, err = ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	ok = (err == nil)
	sub.stopSign = make(chan byte)
	log.Info("recv wait...%v", queueName)
	return err
}

//读取
func (sub *MsgSub) Read(cb rpc.SubFunc) {

	// for d := range sub.readChan {
	// 	log.Info("Received a message: %s", d.Body)
	// 	cb(d.Body, d.Headers)
	// 	d.Ack(false)
	// }
	// return
	for {
		select {
		case d := <-sub.readChan:
			log.Info("Received a message: %s", d.Body)
			cb(d.Body, d.Headers)
			d.Ack(false)
		case <-sub.stopSign:
			sub.onClose()
			break
		}
	}
}

func (sub *MsgSub) Stop() {
	sub.stopSign <- 1
}

func (sub *MsgSub) onClose() {
	if sub.conn != nil {
		sub.conn.Close()
		sub.conn = nil
	}

	if sub.ch != nil {
		sub.ch.Close()
		sub.ch = nil
	}
}
