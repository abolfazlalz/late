package websocket

type Group struct {
	name    string
	clients *Clients
}

func (g *Group) Send(msgType int, msg []byte) error {
	return g.clients.Send(msgType, msg)
}
