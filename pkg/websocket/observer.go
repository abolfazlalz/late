package websocket

type Observer interface {
	Update(ws *Websocket, msg *Message)
	ID() string
}
