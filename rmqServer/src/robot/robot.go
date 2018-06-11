package main

import (
	"base/network"
	"fmt"

	"encoding/json"
	"log"
	"msg"
	"sync"
	"time"
)

type Robot struct {
	account string
	pwd     string

	gateAddr string
	uid      uint64
	key      string

	client network.INetClient
	agent  *Agent
	wg     sync.WaitGroup
}

func NewRobot(account, pwd string) *Robot {
	return &Robot{account: account, pwd: pwd}
}

func (robot *Robot) Start() {
	robot.wg.Add(1)
	if !robot.Login() {
		log.Fatalln("Login fail")
		return
	}
	robot.ConnectGate()
	robot.wg.Wait()
}

func (robot *Robot) Login() bool {
	var addr = *host
	fmt.Println("login...", addr)

	//var addr = "http://47.52.241.13:9900"
	// response, err := http.Get(fmt.Sprintf("%s/login?a=%s&p=1111", addr, robot.account))
	// if err != nil {
	// 	log.Fatalln("login http.get fail:", err)
	// 	return false
	// }
	// defer response.Body.Close()
	// body, _ := ioutil.ReadAll(response.Body)
	// result := msg.C2S_Login{Token: "1"}

	// umErr := gp.Unmarshal(body, &result)
	// if umErr != nil {
	// 	fmt.Println("err:", umErr, "  result:", result)
	// 	return false
	// }
	//fmt.Println("login ok,", result.GateWsAddr)
	// robot.uid = uint64(result.Uid)
	// robot.key = result.Key
	if *nettype == "ws" {
		robot.gateAddr = "ws://" + addr
	} else {
		robot.gateAddr = addr
	}
	return true //result.GetResult() == int32(gameproto.OK)
}

func (robot *Robot) newAgent(conn network.Conn) network.Agent {
	robot.agent = new(Agent)
	robot.agent.conn = conn
	robot.agent.msgHandle = robot.OnMsgRecv
	robot.OnConnected()
	return robot.agent
}

func (robot *Robot) ConnectGate() {
	fmt.Println("ConnectGate:", robot.gateAddr)
	if *nettype == "ws" {
		robot.client = new(network.WSClient)
	} else {
		c := new(network.TCPClient)
		c.LittleEndian = true
		robot.client = c
	}
	robot.client.Set(robot.gateAddr, robot.newAgent)

	//robot.client.LittleEndian = true
	robot.client.Start()

}

func (robot *Robot) OnConnected() {
	fmt.Println("OnConnected...")

	robot.SendMsg("login", &msg.C2S_Login{Token: "1"})
}

func (robot *Robot) EnterGame() {
	fmt.Println("EnterGame...")
	//time.Sleep(2 * time.Second)
	robot.SendMsg("chat", &msg.C2S_Chat{Name: "tom", Msg: "hello~!"})
}

func (robot *Robot) OnMsgRecv(channel byte, msgId interface{}, data []byte) {
	c := 0 //gameproto.ChannelType(channel)
	fmt.Println("OnMsgRecv:", c, " msg:", msgId, " data:", len(data))
	switch msgId {
	case "login":
		robot.EnterGame()
	case "chat":
		var m = new(msg.S2C_Chat)
		json.Unmarshal(data, m)
		fmt.Println("recv chat:", m.Name, m.Msg)
		time.Sleep(5 * time.Second)
		robot.SendMsg("chat", &msg.C2S_Chat{Name: "tom", Msg: "hello~!"})
		//robot.SendMsg("chat", &msg.C2S_Chat{Name: "tom", Msg: "hello~again!"})
	}
}

func (robot *Robot) SendMsg(msgId interface{}, body interface{}) {
	data, err := json.Marshal(body)
	if err != nil {
		fmt.Println("###EncodeMsg error:", err)
		return
	}
	robot.agent.WriteMsg(msgId, data)
}
