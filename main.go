package main

import (
	"asalitermline/pkg/shell"
	"asalitermline/pkg/websocket"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	log.Print("ABL - Dev")
	app := gin.New()

	ws := websocket.NewWebsocket()
	quiteCh := make(chan int)
	go ws.LoopMessage(quiteCh)

	ch := make(chan *shell.Buffer)
	sh := shell.NewShell()
	go func() {
		if err := sh.Execute("ping", "google.com"); err != nil {
			log.Printf("error: %v", err)
		}
	}()

	obs := shell.Observer{}
	ws.Register(obs)

	go func() {
		for {
			msg := <-ch
			log.Print(msg)
		}
	}()

	app.GET("/ws", ws.Handle)
	app.Static("/static", "./public")

	if err := app.Run(":8000"); err != nil {
		log.Printf("error during run app: %v", err)
	}
}
