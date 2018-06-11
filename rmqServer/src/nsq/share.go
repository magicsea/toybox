package main

type AgentMsg struct {
    UID int
    MsgID int
    MsgBody string
}

type PushMsg struct {
	MsgID int
    MsgBody string
}