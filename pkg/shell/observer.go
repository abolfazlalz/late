package shell

import (
	"asalitermline/pkg/websocket"
	websocket2 "github.com/gorilla/websocket"
	"strings"
)

type Observer struct {
}

func (o Observer) Update(ws *websocket.Websocket, msg *websocket.Message) {
	sh := NewShell()
	go func() {
		str := msg.String()
		seps := strings.Split(str, " ")
		var err error
		if len(seps) > 1 {
			err = sh.Execute(seps[0], seps[1:]...)
		} else {
			err = sh.Execute(seps[0])
		}
		if err != nil {
			return
		}
	}()

	go func() {
		for {
			buff := <-sh.bufferCh
			ws.Clients().Send(websocket2.TextMessage, []byte(buff.text))
		}
	}()
}

func (o Observer) ID() string {
	return "shell"
}
