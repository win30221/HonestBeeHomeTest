package main

import (
	"log"
	"math/rand"
	"net"
	"sync"
	"time"
)

const (
	TCP_HOST		= "localhost"
	TCP_PORT		= "9090"
	MAX_CLIENTS		= 50

	CLIENT_PRODUCTION_DELAY	= 200 * time.Millisecond
	CLIENT_SEND_DELAY		= time.Second
)

var (
	WAIT_GROUP		= new(sync.WaitGroup)
	SEND_COMMANDS	= []string{ "Hugo", "HonestBee", "HomeTest" }
)

func main() {

	rand.Seed(time.Now().UnixNano())
	limiter := time.Tick(CLIENT_PRODUCTION_DELAY)
	for id := 1; id <= MAX_CLIENTS; id++ {
		WAIT_GROUP.Add(1)
		go client(id)
		<-limiter
	}
	WAIT_GROUP.Wait()
	log.Println("所有Client皆已退出")
}

// 模擬Client發送訊息直到發送quit退出
func client(id int) {
	addr := TCP_HOST + ":" + TCP_PORT
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println("Client", id, err.Error())
		return
	}
	defer conn.Close()
	limiter := time.Tick(CLIENT_SEND_DELAY)
	for {

		msg := "quit"
		r := rand.Intn(40)

		// 抓取隨機字串
		if r != 0 {
			idx := r % len(SEND_COMMANDS)
			msg = SEND_COMMANDS[idx]
		}

		conn.Write([]byte(msg + "\n"))
		log.Println("Client", id, "發送", msg)

		if msg == "quit" {
			break
		}
		<-limiter
	}
	WAIT_GROUP.Done()
}
