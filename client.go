package main

import "github.com/gorilla/websocket"

type client struct {
	socket *websocket.Conn
	//メッセージが送られるチャンネル
	send chan []byte
	//チャットルーム
	room *room
}

func (c *client) read() {
	defer c.socket.Close()
	//無限ループ
	for {
		if _, msg, err := c.socket.ReadMessage(); err == nil {
			 c.room.forward <- msg
		}else{
			break
		}
	}
}

func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
}
