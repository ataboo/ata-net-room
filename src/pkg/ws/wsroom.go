package ws

import (
	"time"

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

func (r *WSRoom) allPlayers() []*msg.Player {
	players := make([]*msg.Player, len(r.Clients))
	idx := 0
	for _, c := range r.Clients {
		players[idx] = &msg.Player{
			ID:   c.ClientID,
			Name: c.Name,
		}
		idx += 1
	}

	return players
}

func (r *WSRoom) AddUser(conn *websocket.Conn, name string) error {
	r.maxID += 1
	id := r.maxID

	newClient := NewWSClient(conn, id, name)
	r.Clients[id] = newClient

	payload := msg.PlayerIDPayload{
		SubjectID: id,
		PlayerIDs: r.allPlayers(),
	}

	youRes := msg.NewJoinResponse(true, payload)
	if err := WriteResponse(conn, youRes); err != nil {
		common.LogfInfo("failed to send you join %s", err)
	}

	joinRes := msg.NewJoinResponse(false, payload)
	r.BroadcastResponse(&joinRes, id)

	newClient.Start(r.leaveChan, r.reqChan)

	return nil
}

func (r *WSRoom) BroadcastResponse(res *msg.WSResponse, exceptIDs ...int) {
	if res.ID == "" {
		res.ID = msg.GenUniqueID()
	}

	if res.RelayTime == 0 {
		res.RelayTime = time.Now().UnixMilli()
	}

	exceptMap := map[int]bool{}
	for _, id := range exceptIDs {
		exceptMap[id] = true
	}

	for id, c := range r.Clients {
		if _, ok := exceptMap[id]; ok {
			continue
		}

		if !c.WriteResponse(res) {
			r.leaveChan <- c
		}
	}
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

func (r *WSRoom) handleWSRequest(req *msg.WSRequest) {
	switch req.Type {
	case msg.GameEvtReq:
		r.BroadcastResponse(&msg.WSResponse{
			Type:     msg.GameEvtRes,
			SendTime: req.SendTime,
			Sender:   req.Sender,
			ID:       req.ID,
			Name:     req.Name,
			Payload:  req.Payload,
		}, req.Sender)
	case msg.LockReq:
		if r.Locked {
			break
		}
		r.Locked = true
		r.BroadcastResponse(&msg.WSResponse{
			Type:     msg.LockRes,
			SendTime: req.SendTime,
			Sender:   req.Sender,
			ID:       req.ID,
		})
	case msg.UnlockReq:
		if !r.Locked {
			break
		}
		r.Locked = false
		r.BroadcastResponse(&msg.WSResponse{
			Type:     msg.UnlockRes,
			SendTime: req.SendTime,
			Sender:   req.Sender,
			ID:       req.ID,
		})
	default:
		common.LogfInfo("request type not supported: %d", req.Type)
	}
}

func (r *WSRoom) Start() {
	go func() {
		defer func() {
			for _, c := range r.Clients {
				c.conn.Close()
			}

			r.closeChan <- r
		}()

		for {
			select {
			case req := <-r.reqChan:
				common.LogfDebug("Request: %+v", req)
				r.handleWSRequest(req)
			case client := <-r.leaveChan:
				_, ok := r.Clients[client.ClientID]
				if ok {
					common.LogfDebug("removing client: %d", client.ClientID)
					delete(r.Clients, client.ClientID)
					close(client.writeChan)
				}

				if len(r.Clients) == 0 {
					common.LogfDebug("room %s has no clients.", r.Code)
					return
				} else {
					payload := msg.PlayerIDPayload{
						SubjectID: client.ClientID,
						PlayerIDs: r.allPlayers(),
					}
					res := msg.NewLeaveResponse(payload)

					r.BroadcastResponse(&res)
				}
			}
		}
	}()
}
