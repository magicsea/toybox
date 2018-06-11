package gate

import (
	"fmt"
	"net"
	"reflect"

	"base/rpc"
	"base/rpc/rabbitmq"

	"base/log"
	"base/network"
	"base/rpc/proto"
	"base/utils"
	"encoding/json"
	"msg"
	"sync/atomic"
)

type NetType byte

const (
	TCP        NetType = 0
	WEB_SOCKET NetType = 1
)

type Agent interface {
	WriteMsg(msg []byte)
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Close()
	Destroy()
	UserData() interface{}
	SetUserData(data interface{})
	SetDead()
	GetNetType() NetType
}

type GFAgent struct {
	conn     network.Conn
	userData interface{}
	dead     bool
	netType  NetType

	uid       int64
	sessionId string

	mqWrite   rpc.MQClient
	mqRead    *rabbitmq.MsgSubRouter
	readRpcCh chan interface{} //读rpc线程
}

func (a *GFAgent) Active() error {
	a.sessionId = utils.RandomString(10)

	//w
	a.mqWrite = rabbitmq.NewMsgPub()
	err := a.mqWrite.Start("amqp://guest:guest@localhost:5672/", "a2p")
	if err != nil {
		return err
	}

	//r
	a.readRpcCh = make(chan interface{}, 100)
	a.mqRead = rabbitmq.NewMsgSubRouter()
	err = a.mqRead.Start("amqp://guest:guest@localhost:5672/", "p2a", []string{"ss:" + a.sessionId})
	if err != nil {
		return err
	}

	go func() {
		a.mqRead.Read(a.tOnRecvServ)
	}()

	return err
}

func (a *GFAgent) tOnRecvServ(raw []byte, table map[string]interface{}) {
	log.Info("tReadRpc:%v", string(raw))
	if raw == nil {
		panic("bad data")
	}
	m := proto.CSProto{}
	csProtoMk.Decode(raw, &m)
	if m.K == "login" {
		lm := msg.S2C_Login{}
		json.Unmarshal([]byte(m.V), &lm)
		if lm.Result == 0 {
			atomic.StoreInt64(&a.uid, lm.UID)
			a.mqRead.AddRoutKey(fmt.Sprintf("p:%d", lm.UID)) //可以通过uid发消息了
		}
	}
	a.WriteMsg(raw)
}

func (a *GFAgent) GetNetType() NetType {
	return a.netType
}

func (a *GFAgent) SetDead() {
	a.dead = true
}

func (a *GFAgent) Run() {
	for {
		data, err := a.conn.ReadMsg()
		if err != nil {
			log.Error("error on read message: %v", err)
			break
		}
		m := proto.CSProto{}
		csProtoMk.Decode(data, &m)

		uid := atomic.LoadInt64(&a.uid)

		if uid < 1 && m.K != "login" {
			log.Error("no login user,kick:%+v,%v", m, string(data))
			break
		}
		log.Info("agent send:%v", string(data))
		err = a.mqWrite.Pub(data, map[string]interface{}{"UID": fmt.Sprintf("%d", uid), "SID": a.sessionId})
		if err != nil {
			log.Error("error on pub message: %v", err)
			break
		}
	}
}

func (a *GFAgent) OnClose() {
	//todo:not safe
	// if a.agentActor != nil && !a.dead {
	// 	a.agentActor.Tell(&msgs.ClientDisconnect{})
	// }
	if a.mqWrite != nil {
		a.mqWrite.Close()
	}

}

func (a *GFAgent) WriteMsg(data []byte) {
	err := a.conn.WriteMsg(data)
	if err != nil {
		log.Error("write message %v error: %v", reflect.TypeOf(data), err)
	}

}

func (a *GFAgent) LocalAddr() net.Addr {
	return a.conn.LocalAddr()
}

func (a *GFAgent) RemoteAddr() net.Addr {
	return a.conn.RemoteAddr()
}

func (a *GFAgent) Close() {
	a.conn.Close()
}

func (a *GFAgent) Destroy() {
	a.conn.Destroy()

}

func (a *GFAgent) UserData() interface{} {
	return a.userData
}

func (a *GFAgent) SetUserData(data interface{}) {
	a.userData = data
}
