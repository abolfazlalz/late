package websocket

import "encoding/json"

type MessageContent struct {
	Type  string `json:"type"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Message struct {
	sender  *Client
	msgType int
	msg     []byte
	content *MessageContent
}

func NewMessage(sender *Client, msgType int, msg []byte) *Message {
	return &Message{
		sender:  sender,
		msgType: msgType,
		msg:     msg,
		content: nil,
	}
}

func (msg Message) Content() (MessageContent, error) {
	if msg.content != nil {
		return *msg.content, nil
	}
	if err := json.Unmarshal(msg.msg, &msg.content); err != nil {
		return MessageContent{}, err
	}
	return *msg.content, nil
}

func (msg Message) String() string {
	return string(msg.msg)
}

func (msg Message) Sender() *Client {
	return msg.sender
}
