package msg

type C2S_Login struct {
	Token string
}

type S2C_Login struct {
	Result int32
	UID    int64
}

type C2S_Chat struct {
	Name string
	Msg  string
}

type S2C_Chat struct {
	Name string
	Msg  string
}
