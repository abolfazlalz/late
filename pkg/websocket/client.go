package websocket

import (
	"github.com/gorilla/websocket"
	"github.com/samber/lo"
	"log"
)

type Client struct {
	id   int
	conn *websocket.Conn
}

func (c *Client) Send(msgType int, msg []byte) error {
	return c.conn.WriteMessage(msgType, msg)
}

func (c *Client) LoopMessage(msgChan chan *Message) {
	for {
		msgType, msg, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("error durig read message: %v", err)
			return
		}

		msgChan <- NewMessage(c, msgType, msg)
	}
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
		id:   c.lastId,
		conn: conn,
	}
	c.clients = append(c.clients, client)

	return client
}

func (c *Clients) Delete(id int) {
	c.clients = lo.Filter(c.clients, func(item *Client, _ int) bool {
		return item.id != id
	})
}

func (c *Clients) Send(msgType int, msg []byte) error {
	for _, client := range c.clients {
		if err := client.Send(msgType, msg); err != nil {
			return err
		}
	}
	return nil
}
