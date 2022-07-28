package ws

import (
	"github.com/ataboo/ata-net-room/pkg/ws/msg"
	"github.com/gorilla/websocket"
)

type WSRoom struct {
	Code      string
	Capacity  int
	Locked    bool
	Clients   map[int]*WSClient
	MaxID     int
	leaveChan chan *WSClient
	reqChan   chan *msg.WSRequest
}

func (r *WSRoom) AddUser(conn *websocket.Conn) error {
	r.MaxID += 1
	id := r.MaxID

	otherPlayerIDs := make([]int, len(r.Clients))
	for i, c := range r.Clients {
		otherPlayerIDs[i] = c.ClientID
	}

	payload := msg.PlayerIDPayload{
		PlayerID:        id,
		OtherPlayersIDs: otherPlayerIDs,
	}

	youRes := msg.NewJoinResponse(true, payload)
	if err := WriteResponse(conn, youRes); err != nil {
		return err
	}

	joinRes := msg.NewJoinResponse(false, payload)
	r.BroadcastResponse(joinRes)

	r.Clients[id] = NewWSClient(conn, id)
	r.Clients[id].Start(r.leaveChan, r.reqChan)

	return nil
}

func (r *WSRoom) BroadcastResponse(res msg.WSResponse) {

}

func NewWSRoom(req *msg.WSJoinRequest) *WSRoom {
	return &WSRoom{
		Code:      req.RoomCode,
		Capacity:  req.RoomSize,
		Locked:    false,
		Clients:   map[int]*WSClient{},
		MaxID:     0,
		leaveChan: make(chan *WSClient),
		reqChan:   make(chan *msg.WSRequest),
	}
}

func (r *WSRoom) Start() {
	go func() {
		select {
		case req := <-r.reqChan:
			break
		case client := <-r.leaveChan:
			_, ok := r.Clients[client.ClientID]
			if ok {
				delete(r.Clients, client.ClientID)
				close(client.writeChan)
			}

			break
		}
	}()
}
