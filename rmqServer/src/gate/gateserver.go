package gate

import (
	"base/log"
	"base/network"
	"base/rpc/proto"
	"encoding/json"
	"msg"
	_ "net/http/pprof"
)

var csProtoMk *proto.JsonProtoMaker

func Run() {
	log.Info("gate run...")

	csProtoMk = new(proto.JsonProtoMaker)

	tcpAddr := "127.0.0.1:7200"
	var tcpServer *network.TCPServer
	tcpServer = new(network.TCPServer)
	tcpServer.Addr = tcpAddr
	tcpServer.MaxConnNum = 10000
	tcpServer.PendingWriteNum = 1000
	tcpServer.LenMsgLen = 2
	tcpServer.MaxMsgLen = 50000
	tcpServer.LittleEndian = true
	tcpServer.NewAgent = func(conn *network.TCPConn) network.Agent {
		a := &GFAgent{conn: conn, netType: TCP}
		err := a.Active()
		if err != nil {
			log.Error("NewAgent fail:%v", err.Error())
		}
		return a
	}

	if tcpServer != nil {
		tcpServer.Start()
	}

	select {}

	log.Info("gate end!")
}

func TestSend() {
	a := new(GFAgent)
	a.Active()

	j, _ := json.Marshal(&msg.C2S_Login{Token: "1"})
	m := proto.CSProto{K: "login", V: string(j)}
	bt, _ := json.Marshal(&m)
	a.mqWrite.Pub(bt, map[string]interface{}{"UID": "0", "SID": a.sessionId})
	select {}
}
