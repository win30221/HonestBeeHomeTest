package server

import (
	"log"
	"net"
	"strings"
	"time"
)

type Server struct {
	listener	net.Listener	// 監聽連線
	clients		map[int]*Client	// Client列表
	joinsniffer	chan net.Conn	// Client連線
	insniffer	chan string		// Client發訊
	quitsniffer chan *Client	// Client退出
}

type Client struct {
	id		int				// Client ID
	conn	net.Conn		// Client連線
	in		chan string		// Client收訊
	quit	chan *Client	// Client退出
}

func (this *Server) joinHandler(conn net.Conn) {
	// 建立新的Client
	client := CreateClient(CLIENT_ID, conn)
	this.clients[client.id] = client
	log.Printf("Client[%d]連線", client.id)
	CLIENT_ID++
	// 持續監聽Client發訊
	go func() {
		for {
			msg := <-client.in
			this.insniffer <- msg
		}
	}()
	// 監聽Client退出
	go func() {
		quit := <- client.quit
		this.quitsniffer <- quit
	}()
}

func (this *Server) receivedHandler(message string) {
	log.Println("收到的訊息:", message)
}

func (this *Server) quitHandler(client *Client) {
	if client != nil {
		delete(this.clients, client.id)
		log.Printf("Client[%d]退出", client.id)
		client.conn.Close()
	}
}

func CreateClient(id int, conn net.Conn) *Client {
	client := &Client{
		id:   id,
		conn: conn,
		in:   make(chan string),
		quit: make(chan *Client),
	}
	go client.read()
	return client
}

// 實際讀取Client訊息
func (this *Client) read() {
	for {
		data := make([]byte, 1024)
		n, err := this.conn.Read(data)
		if err != nil {
			this.quiting()
			return
		}
		// 去除尾部換行塞入in Channel
		this.in <- strings.TrimSuffix(string(data[0:n]), "\n")
	}
}

// Client呼叫退出
func (this *Client) quiting() {
	this.quit <- this
}

// 監聽所有Client與Server溝通的Channel
func (this *Server) listen() {
	go func() {
		limiter := time.Tick(time.Second / time.Duration(REQUEST_PER_SECOND))
		for {
			select {
			// 接收到了訊息
			case message := <-this.insniffer:
				<-limiter
				go this.receivedHandler(message)
			// 新來了一個連線
			case conn := <-this.joinsniffer:
				go this.joinHandler(conn)
			// 退出了一個連線
			case client := <-this.quitsniffer:
				go this.quitHandler(client)
			}
		}
	}()
}

// 監聽Client建立連線
func (this *Server) start() {
	addr := "0.0.0.0:" + TCP_PORT
	this.listener, _ = net.Listen("tcp", addr)
	log.Printf("開始監聽TCP Server [%s] Port", TCP_PORT)
	for {
		conn, err := this.listener.Accept()
		if err != nil {
			log.Fatalln(err)
			return
		}
		// 將連線給Channel統一管理
		this.joinsniffer <- conn
	}
}

func StartTCPServer() {
	server := &Server{
		clients:		make(map[int]*Client, MAX_CLIENTS),
		joinsniffer:	make(chan net.Conn),
		quitsniffer:	make(chan *Client),
		insniffer:		make(chan string, MAX_REQUEST_COUNT),
	}
	server.listen()
	server.start()
}