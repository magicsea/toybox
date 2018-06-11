package game

import (
	"base/rpc/proto"
	"encoding/json"
)

type Content struct {
	msg proto.CSProto
	uid int64
	sid string
}

func (c *Content) GetUID() int64 {
	return c.uid
}

func (c *Content) GetMsgID() string {
	return c.msg.K
}

func (c *Content) GetBody() string {
	return c.msg.V
}

func (c *Content) GetBodyObj(obj interface{}) error {
	err := json.Unmarshal([]byte(c.GetBody()), obj)
	return err
}

//返回消息
func (c *Content) Write(msgID string, bodyObj interface{}) {
	SendPlayer(c.uid, msgID, bodyObj)
}
