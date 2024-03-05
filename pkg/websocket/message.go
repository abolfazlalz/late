package websocket

type Message struct {
	sender  *Client
	msgType int
	msg     []byte
}

func NewMessage(sender *Client, msgType int, msg []byte) *Message {
	return &Message{
		sender:  sender,
		msgType: msgType,
		msg:     msg,
	}
}

func (msg Message) String() string {
	return string(msg.msg)
}

func (msg Message) Sender() *Client {
	return msg.sender
}
