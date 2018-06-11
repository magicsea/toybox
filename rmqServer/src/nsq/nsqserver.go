package main

import (
	"encoding/json"
	"fmt"
	"github.com/nsqio/go-nsq"
	"sync"
)

type NSQHandler struct {
	conId int
}

func (this *NSQHandler) HandleMessage(msg *nsq.Message) error {
	fmt.Println("receive:", this.conId, "message:", string(msg.Body))
	amsg := AgentMsg{}
	json.Unmarshal(msg.Body, &amsg)
	nsqd := "magicsea.top:4150"
	producer, _ := nsq.NewProducer(nsqd, nsq.NewConfig())
	topic := fmt.Sprintf("push%d", amsg.UID)
	pmsg := PushMsg{MsgID: 1, MsgBody: "rep," + amsg.MsgBody}
	pdata, _ := json.Marshal(&pmsg)
	producer.Publish(topic, pdata)
	return nil
}

func testNSQ() {
	waiter := sync.WaitGroup{}
	waiter.Add(1)

	go func() {
		defer waiter.Done()
		config := nsq.NewConfig()
		config.MaxInFlight = 9

		//建立多个连接
		for i := 0; i < 2; i++ {
			consumer, err := nsq.NewConsumer("a2s", "default", config)
			if nil != err {
				fmt.Println("err", err)
				return
			}

			consumer.AddHandler(&NSQHandler{conId: i})
			err = consumer.ConnectToNSQD("magicsea.top:4150")
			if nil != err {
				fmt.Println("err", err)
				return
			}
		}
		select {}

	}()

	waiter.Wait()
}
func main() {
	testNSQ()

}
