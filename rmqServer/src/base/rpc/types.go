package rpc

type MsgPack interface {
	GetKey() interface{}
	GetBody() interface{}
}

type ProtoMaker interface {
	Encode(pack MsgPack) ([]byte, error)
	Decode(raw []byte, out MsgPack) error
}

type MQClient interface {
	Start(addr string, queueName string) error
	Close()
	Pub(body []byte, table map[string]interface{}) error
}

type SubFunc func(data []byte, table map[string]interface{})

type MQServer interface {
	Start(addr string, queueName string, routerKeys []string) error
	Read(cb SubFunc)
	Stop()
}
