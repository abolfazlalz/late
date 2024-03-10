package shell

import (
	"asalitermline/pkg/websocket"
	gorillaws "github.com/gorilla/websocket"
	"log"
	"strings"
	"sync"
)

type Observer struct {
	id      int
	execute *CommandExecute
	mu      sync.Mutex
	msgCh   chan *Buffer
}

func NewObserver() *Observer {
	return &Observer{
		msgCh: make(chan *Buffer),
		mu:    sync.Mutex{},
	}
}

func (o *Observer) Update(ws *websocket.Websocket, msg *websocket.Message) {
	sh := NewShell()

	o.mu.Lock()
	defer o.mu.Unlock()

	content, _ := msg.Content()
	str := content.Value
	seps := strings.Split(str, " ")
	if len(seps) > 1 {
		o.id, o.execute = sh.Execute(seps[0], seps[1:]...)
	} else {
		o.id, o.execute = sh.Execute(seps[0])
	}

	cmdQuit := make(chan int)
	msg.Sender().OnClose(cmdQuit)
	o.execute.Subscribe(o.msgCh)

	go func() {
		err := o.execute.Run(cmdQuit)
		if err != nil {
			log.Printf("error during run Command: %v", err)
			return
		}
	}()

	go o.run(ws, msg)
}

func (o *Observer) run(_ *websocket.Websocket, msg *websocket.Message) {
	for {
		buff := <-o.msgCh
		err := msg.Sender().Send(gorillaws.TextMessage, []byte(buff.text))
		if err != nil {
			log.Print("error during send message to client:", err)
			continue
		}
	}
}

func (o *Observer) ID() string {
	return "shell"
}
