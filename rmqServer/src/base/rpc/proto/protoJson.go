package proto

import (
	"base/log"
	"base/rpc"
	"encoding/json"
)

type MsgPackJson struct {
	uid  int64
	key  string
	body string
}

func (p *MsgPackJson) GetKey() interface{} {
	return p.key
}
func (p *MsgPackJson) GetBody() interface{} {
	return p.body
}
func (p *MsgPackJson) GetUID() int64 {
	return p.uid
}

//=============================================================================
type JsonProtoMaker struct {
}

func (m *JsonProtoMaker) Encode(pack rpc.MsgPack) ([]byte, error) {
	data, err := json.Marshal(pack)
	log.PrintOnError(err, "json.Marshal")
	return data, err
}

func (m *JsonProtoMaker) Decode(raw []byte, out interface{}) error {
	err := json.Unmarshal(raw, out)
	log.PrintOnError(err, "json.Unmarshal:"+string(raw))
	return err
}
