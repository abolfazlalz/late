package websocket

import (
	"github.com/gorilla/websocket"
	"github.com/samber/lo"
	"log"
)

type ClientStatus int

const (
	ClientConnected ClientStatus = iota
	ClientDisconnected
)

type Client struct {
	id     int
	conn   *websocket.Conn
	status ClientStatus
	quitCh []chan int
}

func (c *Client) Send(msgType int, msg []byte) error {
	if c.status == ClientDisconnected {
		return nil
	}
	return c.conn.WriteMessage(msgType, msg)
}

func (c *Client) LoopMessage(msgChan chan *Message) {
	for {
		msgType, msg, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("error durig read message: %v", err)
			c.quit()
			c.status = ClientDisconnected
			return
		}

		msgChan <- NewMessage(c, msgType, msg)
	}
}

func (c *Client) quit() {
	for _, quit := range c.quitCh {
		quit <- 1
	}
}

func (c *Client) OnClose(quitCh chan int) {
	c.quitCh = append(c.quitCh, quitCh)
}

type Clients struct {
	lastId  int
	clients []*Client
}

func NewClients() *Clients {
	return &Clients{clients: make([]*Client, 0), lastId: 0}
}

func (c *Clients) AddConn(conn *websocket.Conn) *Client {
	c.lastId++
	client := &Client{
		id:     c.lastId,
		conn:   conn,
		status: ClientConnected,
		quitCh: make([]chan int, 0),
	}
	c.clients = append(c.clients, client)

	return client
}

func (c *Clients) Delete(id int) {
	c.clients = lo.Filter(c.clients, func(item *Client, _ int) bool {
		return item.id != id
	})
}

func (c *Clients) Quit() {
	for _, client := range c.clients {
		for _, ch := range client.quitCh {
			ch <- 1
		}
	}
}

func (c *Clients) Send(msgType int, msg []byte) error {
	for _, client := range c.clients {
		if client.status == ClientDisconnected {
			continue
		}
		if err := client.Send(msgType, msg); err != nil {
			return err
		}
	}
	return nil
}
