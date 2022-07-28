package ws

import (
	"time"

	"github.com/ataboo/ata-net-room/pkg/common"
	"github.com/ataboo/ata-net-room/pkg/ws/msg"
	"github.com/gorilla/websocket"
)

type WSClient struct {
	ClientID  int
	conn      *websocket.Conn
	writeChan chan *msg.WSResponse
}

func NewWSClient(conn *websocket.Conn, id int) *WSClient {
	return &WSClient{
		conn:      conn,
		writeChan: make(chan *msg.WSResponse),
		ClientID:  id,
	}
}

func (c *WSClient) Start(leaveChan chan<- *WSClient, reqChan chan<- *msg.WSRequest) {
	go c.readPump(leaveChan, reqChan)
	go c.writePump(leaveChan)
}

func (c *WSClient) readPump(leaveChan chan<- *WSClient, reqChan chan<- *msg.WSRequest) {
	defer func() {
		leaveChan <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(MaxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(PongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(PongWait)); return nil })
	for {
		c.conn.SetReadDeadline(time.Now().Add(ReadWait))
		req := msg.WSRequest{}
		err := c.conn.ReadJSON(req)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				common.LogfInfo("error: %v", err)
			}
			return
		}

		reqChan <- &req
	}
}

func (c *WSClient) writePump(leaveChan chan<- *WSClient) {
	ticker := time.NewTicker(PingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(WriteWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		case res, ok := <-c.writeChan:
			if !ok {
				c.conn.SetWriteDeadline(time.Now().Add(WriteWait))
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.conn.SetWriteDeadline(time.Now().Add(WriteWait))
			if err := c.conn.WriteJSON(res); err != nil {
				return
			}
		}
	}
}

func (c *WSClient) WriteResponse(res *msg.WSResponse) bool {
	select {
	case c.writeChan <- res:
		return true
	default:
		close(c.writeChan)
		return false

	}
}

func WriteResponse(conn *websocket.Conn, res msg.WSResponse) error {
	conn.SetWriteDeadline(time.Now().Add(WriteWait))
	return conn.WriteJSON(res)
}
