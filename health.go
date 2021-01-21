package main

import (
	"fmt"
	"github.com/zinvapel/timetracker/contract"
	"log"
	"net"
	"time"
)

func health() {
	conn, err := net.Listen("tcp", ":50000")
	if err != nil {
		log.Println("[health] Server start failed", err)
		return
	}
	defer conn.Close()

	for {
		cn, err := conn.Accept()

		if err != nil {
			log.Println("[health] Unable to connect with client")
			continue
		}

		go func() {
			defer cn.Close()
			ch := make(chan bool)
			contract.GetEventBus().Push(contract.NewEvent("health", ch))

			select {
			case <-ch:
				_, err := fmt.Fprintln(cn, "ok")
				if err != nil {
					log.Printf("[health] failed %v", err)
				}
			case <-time.After(time.Second):
				log.Println("[health] failed by timeout")
				return
			}
		}()
	}
}
