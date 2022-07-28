package ws

import "github.com/gorilla/websocket"

type WSClient struct {
	Conn     websocket.Conn
	ClientID int
}

func (c *WSClient) Start(stopChan chan bool) {

}
