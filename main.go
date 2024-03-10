package main

import (
	"asalitermline/pkg/shell"
	"asalitermline/pkg/websocket"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
)

func main() {
	log.Print("ABL - Dev")
	app := gin.New()

	signalC := make(chan os.Signal, 1)
	signal.Notify(signalC, os.Interrupt)

	ws := websocket.NewWebsocket()
	quiteCh := make(chan int)
	go ws.LoopMessage(quiteCh)

	obs := shell.NewObserver()
	ws.Register(obs)

	app.GET("/ws", ws.Handle)
	app.DELETE("kill/:id", func(c *gin.Context) {
		sh := shell.NewShell()
		idStr := c.Param("id")
		id, _ := strconv.Atoi(idStr)
		if err := sh.Kill(id); err != nil {
			log.Printf("error during kill session: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "error during kill session"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"msg": "ok"})
	})
	app.GET("/ls", func(c *gin.Context) {
		sh := shell.NewShell()
		c.JSON(http.StatusOK, sh.List())
	})
	app.Static("/static", "./resources")

	if err := app.Run(":8000"); err != nil {
		log.Printf("error during run app: %v", err)
	}
}
