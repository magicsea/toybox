package main

import (
	"flag"
	"fmt"
)

var acc = flag.String("u", "magicse_1", "account")
var host = flag.String("host", "127.0.0.1:7200", "login url")
var nettype = flag.String("net", "tcp", "net type:tcp or ws")

func main() {
	flag.Parse()
	fmt.Println("start:", *acc, *host, *nettype)
	r := NewRobot(*acc, "111")
	r.Start()
	fmt.Println("end")
}
