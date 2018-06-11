package main

import (
	"base/network"
	//	"io/ioutil"
	"log"
	//	"net/http"
	//	"github.com/magicsea/ganet/network/protobuf"
	//	"github.com/gogo/protobuf/proto"
	"base/rpc/proto"
	"encoding/json"
)

func newAgent(conn network.Conn) network.Agent {
	Client := new(Agent)
	Client.conn = conn
	return Client
}

type Agent struct {
	conn      network.Conn
	msgHandle func(channel byte, msgId interface{}, data []byte)
}

func (a *Agent) Run() {
	log.Println("Agent.run")
	for {
		data, err := a.conn.ReadMsg()
		if err != nil {
			log.Println("read message: ", err)
			break
		}
		var jd = new(proto.CSProto)
		errUm := json.Unmarshal(data, jd)
		if errUm != nil {
			log.Println("Unmarshal fail:", errUm)
			break
		}
		a.msgHandle(0, jd.K, []byte(jd.V))
	}
}

func (a *Agent) OnClose() {}

/*
func (a *Agent) WriteMsg(channel byte, msgId byte, msg []byte) {

	data := []byte{channel, msgId}
	data = append(data, msg...)
	err := a.conn.WriteMsg(data)
	if err != nil {
		log.Println("write message error:", err)
	}

}
*/
func (a *Agent) WriteMsg(msgID interface{}, rawmsg []byte) {

	var jd = proto.CSProto{K: msgID.(string), V: string(rawmsg)}
	//m := map[string]interface{}{JsonIdName: msgID,JsonMsgName:rawmsg}
	data, _ := json.Marshal(jd)

	// data := []byte{msgId, 0, 0}
	// data = append(data, msg...)
	serr := a.conn.WriteMsg(data)
	if serr != nil {
		log.Println("write message error:", serr)
	}

}

//func (a *Agent) LocalAddr() net.Addr {
//return a.conn.LocalAddr()
//}

//func (a *Agent) RemoteAddr() net.Addr {
//	return a.conn.RemoteAddr()
//}

func (a *Agent) Close() {
	a.conn.Close()
}

func (a *Agent) Destroy() {
	a.conn.Destroy()
}
