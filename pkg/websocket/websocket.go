package websocket

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type Websocket struct {
	ids       int
	upgrader  websocket.Upgrader
	clients   *Clients
	msgChan   chan *Message
	Observers []Observer
}

func NewWebsocket() *Websocket {
	return &Websocket{
		ids: 0,
		upgrader: websocket.Upgrader{
			WriteBufferSize: 1024,
			ReadBufferSize:  1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		clients: NewClients(),
		msgChan: make(chan *Message),
	}
}

func (ws *Websocket) Handle(c *gin.Context) {
	conn, err := ws.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("error during ws connection: %v", err)
	}
	client := ws.clients.AddConn(conn)
	client.LoopMessage(ws.msgChan)
}

func (ws *Websocket) Register(observer Observer) {
	ws.Observers = append(ws.Observers, observer)
}

func (ws *Websocket) Clients() *Clients {
	return ws.clients
}

func (ws *Websocket) notify(msg *Message) {
	content, err := msg.Content()
	if err != nil {
		return
	}
	for _, observer := range ws.Observers {
		if observer.ID() == content.Type {
			observer.Update(ws, msg)
		}
	}
}

func (ws *Websocket) LoopMessage(quiteCh chan int) {
	for {
		select {
		case <-quiteCh:
			ws.clients.Quit()
			return
		case msg := <-ws.msgChan:
			log.Print(msg)
			ws.notify(msg)
			break
		}
	}
}
