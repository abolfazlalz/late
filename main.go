package main

import (
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

	app.GET("/ws", ws.Handle)

	if err := app.Run(":8000"); err != nil {
		log.Printf("error during run app: %v", err)
	}
}
