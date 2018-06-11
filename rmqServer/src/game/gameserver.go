package game

import (
	"base/log"
	"base/rpc/proto"
	"base/rpc/rabbitmq"
	"base/utils"
	"encoding/json"
	"fmt"
	"strconv"
)

func onRecv(raw []byte, table map[string]interface{}) {
	defer utils.PrintPanicStack()
	log.Info("onRecv:%v,%v,%v", table["SID"], table["UID"], string(raw))
	m := proto.CSProto{}
	csProtoMk.Decode(raw, &m)
	sid := table["SID"].(string)
	uid, _ := strconv.Atoi(table["UID"].(string))
	go router(m, int64(uid), sid)
}

func router(msg proto.CSProto, uid int64, sid string) {
	c := &Content{msg: msg, uid: uid, sid: sid}
	switch msg.K {
	case "login":
		OnLogin(c)
	case "chat":
		OnChat(c)
	}
}

var csProtoMk *proto.JsonProtoMaker
var pushAgent *rabbitmq.MsgPubRouter //write

func Run() {

	csProtoMk = new(proto.JsonProtoMaker)
	addr := "amqp://guest:guest@localhost:5672/"
	//c
	pushAgent = new(rabbitmq.MsgPubRouter)
	err := pushAgent.Start(addr, "p2a")
	if err != nil {
		log.Error("%v", err)
		return
	}
	//s
	var mqServ = rabbitmq.NewMsgSub() //read
	err = mqServ.Start(addr, "a2p")
	if err != nil {
		log.Error("%v", err)
		return
	}

	log.Info("game run...")
	mqServ.Read(onRecv)
	log.Info("game end:%v", err)
}

func SendPlayer(uid int64, msgId string, obj interface{}) {
	data, _ := json.Marshal(obj)
	m := proto.CSProto{K: msgId, V: string(data)}
	raw, _ := json.Marshal(&m)
	pushAgent.Pub(raw, nil, fmt.Sprintf("p:%d", uid))
}

func SendAgent(sid string, msgId string, obj interface{}) {
	data, _ := json.Marshal(obj)
	m := proto.CSProto{K: msgId, V: string(data)}
	raw, _ := json.Marshal(&m)
	pushAgent.Pub(raw, nil, "ss:"+sid)
}
