package main

import "github.com/win30221/HonestBeeHomeTest/server"

func main() {
	go func() {
		server.StartHttpServer()
	}()
	server.StartTCPServer()
}
