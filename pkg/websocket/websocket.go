package websocket

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
)

type Websocket struct {
	ids      int
	upgrader websocket.Upgrader
	clients  *Clients
	msgChan  chan *Message
}

func NewWebsocket() *Websocket {
	return &Websocket{
		ids: 0,
		upgrader: websocket.Upgrader{
			WriteBufferSize: 1024,
			ReadBufferSize:  1024,
		},
		clients: NewClients(),
		msgChan: make(chan *Message),
	}
}

func (ws Websocket) Handle(c *gin.Context) {
	conn, err := ws.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("error during ws connection: %v", err)
	}
	client := ws.clients.AddConn(conn)
	client.LoopMessage(ws.msgChan)
}

func (ws Websocket) LoopMessage(quiteCh chan int) {
	for {
		select {
		case <-quiteCh:
			return
		case msg := <-ws.msgChan:
			log.Print(msg)
			break
		}
	}
}
