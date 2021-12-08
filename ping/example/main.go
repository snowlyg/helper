package main

import (
	"fmt"

	"github.com/snowlyg/helper/ping"
)

func main() {
	ip := "10.0.0.113"
	ok, msg := ping.GetPingMsg(ip)
	if !ok {
		fmt.Printf("%s ping is fault,get msg %s \n", ip, msg)
	}
}
