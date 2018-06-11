package main

import (
	"encoding/json"
	"fmt"
	"github.com/nsqio/go-nsq"
	"time"
)

var producer *nsq.Producer

func main() {
	id := 2
	initListen(id)

	nsqd := "magicsea.top:4150"
	producer, err := nsq.NewProducer(nsqd, nsq.NewConfig())
	for i := 0; i < 9; i++ {
		msg := fmt.Sprintf("nihao%d", i)
		amsg := AgentMsg{UID: id, MsgID: i, MsgBody: msg}
		data, _ := json.Marshal(&amsg)
		producer.Publish("a2s", []byte(data))
		fmt.Println("send ", i)
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Second)

	}
	producer.Stop()
	select {}
}

func initListen(id int) {
	topic := fmt.Sprintf("push%d", id)
    config := nsq.NewConfig()
    nsq.NewConn().
	consumer, err := nsq.NewConsumer(topic, "default", config)
	if nil != err {
		fmt.Println("err", err)
		return
	}

	consumer.AddHandler(&NSQListenHandler{conId: id})
	err = consumer.ConnectToNSQD("magicsea.top:4150")
	if nil != err {
		fmt.Println("err", err)
		return
	}
}

type NSQListenHandler struct {
	conId int
}

func (this *NSQListenHandler) HandleMessage(msg *nsq.Message) error {
	fmt.Println("listen:", this.conId, "message:", string(msg.Body))

	return nil
}
