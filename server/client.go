package server

import (
	"net"
	"strings"
)

type Client struct {
	id		int				// Client ID
	conn	net.Conn		// Client連線
	in		chan string		// Client收訊
	quit	chan *Client	// Client退出
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
	close(this.in)
	this.quit <- this
}