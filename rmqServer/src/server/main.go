package main

import (
	"base/log"
	"game"
	"gate"
)

func main() {

	err := log.NewLogGroup("debug", "log", true, 7)
	if err != nil {
		panic(err)
	}
	defer log.Close()

	go game.Run()
	go gate.Run()
	//gate.TestSend()
	select {}
}
