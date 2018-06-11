package rabbitmq

import (
	"base/log"
	"base/rpc"
	"github.com/streadway/amqp"
)

type MsgSubRouter struct {
	ch        *amqp.Channel
	conn      *amqp.Connection
	queueName string
	exName    string

	readChan <-chan amqp.Delivery
	stopSign chan byte
}

func NewMsgSubRouter() *MsgSubRouter {
	return new(MsgSubRouter)
}

func (sub *MsgSubRouter) Start(addr string, exName string, routerKeys []string) error {
	var ok bool
	//异常退出释放
	defer func() {
		if !ok {
			//sub.onClose()
		}
	}()

	conn, err := amqp.Dial(addr)
	if err != nil {
		return err
	}
	//defer conn.Close()
	sub.conn = conn

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	//defer ch.Close()
	sub.ch = ch

	err = ch.ExchangeDeclare(
		exName,   // name
		"direct", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return err
	}

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	// err = ch.Qos(
	// 	1,     // prefetch count
	// 	0,     // prefetch size
	// 	false, // global
	// )
	// if err != nil {
	// 	return err
	// }

	for _, routerKey := range routerKeys {
		// log.Printf("Binding queue %s to exchange %s with routing key %s",
		// 	q.Name, "logs_direct", s)
		err = ch.QueueBind(
			q.Name,    // queue name
			routerKey, // routing key
			exName,    // exchange
			false,
			nil)
		log.PrintOnError(err, "Failed to bind a queue")
	}
	sub.queueName = q.Name
	sub.exName = exName

	sub.readChan, err = ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	log.PrintOnError(err, "Failed to register a consumer")
	log.Info("recv wait...%v:%+v", exName, routerKeys)
	sub.stopSign = make(chan byte)
	ok = (err == nil)
	return err
}

func (sub *MsgSubRouter) AddRoutKey(key string) error {
	err := sub.ch.QueueBind(
		sub.queueName, // queue name
		key,           // routing key
		sub.exName,    // exchange
		false,
		nil)
	log.PrintOnError(err, "Failed to bind a queue")
	return err
}

//读取
func (sub *MsgSubRouter) Read(cb rpc.SubFunc) {
	for {
		select {
		case d := <-sub.readChan:
			log.Info("Received a message: %s", d.Body)
			cb(d.Body, d.Headers)
			//d.Ack(false)
		case <-sub.stopSign:
			sub.onClose()
			break
		}
	}
}

func (sub *MsgSubRouter) Stop() {
	sub.stopSign <- 1
}

func (sub *MsgSubRouter) onClose() {
	log.Info("MsgSubRouter close")
	if sub.conn != nil {
		sub.conn.Close()
		sub.conn = nil
	}

	if sub.ch != nil {
		sub.ch.Close()
		sub.ch = nil
	}
}
