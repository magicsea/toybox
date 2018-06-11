package game

import (
	"base/log"
	"msg"
	"strconv"
)

func OnLogin(c *Content) {
	log.Info("OnLogin:%d=>%v", c.GetUID(), c.GetBody())
	m := new(msg.C2S_Login)
	c.GetBodyObj(m)
	id, _ := strconv.Atoi(m.Token)
	SendAgent(c.sid, "login", &msg.S2C_Login{Result: 0, UID: int64(id)})
}

func OnChat(c *Content) {
	log.Info("OnChat:%d=>%v", c.GetUID(), c.GetBody())
	m := new(msg.C2S_Chat)
	c.GetBodyObj(m)
	c.Write("chat", &msg.S2C_Chat{Name: m.Name, Msg: m.Msg})
}
