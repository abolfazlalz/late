package websocket

import "fmt"

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
	return fmt.Sprintf("received message from %d -> type: %d text: %s", msg.sender.id, msg.msgType, msg.msg)
}
