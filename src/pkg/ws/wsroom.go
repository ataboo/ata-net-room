package ws

import (
	"github.com/ataboo/ata-net-room/pkg/common"
	"github.com/ataboo/ata-net-room/pkg/ws/msg"
	"github.com/gorilla/websocket"
)

type WSRoom struct {
	Code      string
	Capacity  int
	Locked    bool
	Clients   map[int]*WSClient
	maxID     int
	leaveChan chan *WSClient
	reqChan   chan *msg.WSRequest
	closeChan chan<- *WSRoom
}

func (r *WSRoom) AddUser(conn *websocket.Conn) error {
	r.maxID += 1
	id := r.maxID

	otherPlayerIDs := make([]int, len(r.Clients))
	for i, c := range r.Clients {
		otherPlayerIDs[i] = c.ClientID
	}

	payload := msg.PlayerIDPayload{
		SubjectID: id,
		PlayerIDs: otherPlayerIDs,
	}

	youRes := msg.NewJoinResponse(true, payload)
	if err := WriteResponse(conn, youRes); err != nil {
		common.LogfInfo("failed to send you join %s", err)
	}

	joinRes := msg.NewJoinResponse(false, payload)
	r.BroadcastResponse(joinRes)

	r.Clients[id] = NewWSClient(conn, id)
	r.Clients[id].Start(r.leaveChan, r.reqChan)

	return nil
}

func (r *WSRoom) BroadcastResponse(res msg.WSResponse) {
	panic("not implemented")
}

func NewWSRoom(req *msg.WSJoinRequest, closeChan chan<- *WSRoom) *WSRoom {
	return &WSRoom{
		Code:      req.RoomCode,
		Capacity:  req.RoomSize,
		Locked:    false,
		Clients:   map[int]*WSClient{},
		maxID:     0,
		leaveChan: make(chan *WSClient),
		reqChan:   make(chan *msg.WSRequest),
		closeChan: closeChan,
	}
}

func (r *WSRoom) Start() {
	go func() {
		select {
		case req := <-r.reqChan:
			common.LogfDebug("%+v", req)
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
