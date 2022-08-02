package ws

import (
	"encoding/json"
	"net/http"
	"regexp"
	"time"

	"github.com/ataboo/ata-net-room/pkg/common"
	"github.com/ataboo/ata-net-room/pkg/ws/msg"
	"github.com/gorilla/websocket"
)

const (
	MaxMessageSize = 2048
	ReadWait       = 3 * time.Second
	WriteWait      = 3 * time.Second
	PongWait       = 10 * time.Second
	PingPeriod     = 5 * time.Second
)

type ServerConfig struct {
	Host         string
	RoomCapacity int
	Subprotocol  string
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
		common.LogfDebug("WS Error: %d, %s\n", status, reason.Error())
	},
	Subprotocols: []string{"atanet_v1"},
}

type WSService struct {
	config        ServerConfig
	Rooms         map[string]*WSRoom
	roomCloseChan chan *WSRoom
}

func (s *WSService) HandleJoin(w http.ResponseWriter, r *http.Request) {
	if !s.hasSubProtocolHeader(r) {
		common.LogfDebug("protocol not supported")
		http.Error(w, "invalid protocol", http.StatusBadRequest)
		return
	}

	h := http.Header{}
	h.Set("Sec-Websocket-Protocol", s.config.Subprotocol)

	conn, err := upgrader.Upgrade(w, r, h)
	if err != nil {
		common.LogfDebug("failed to upgrade: %s", err.Error())
		http.Error(w, "failed to upgrade", http.StatusUnprocessableEntity)
		return
	}

	conn.SetReadLimit(MaxMessageSize)
	conn.SetReadDeadline(time.Now().Add(ReadWait))
	mType, msgBytes, err := conn.ReadMessage()
	if mType != websocket.BinaryMessage || err != nil {
		common.LogfDebug("expected binary message")
		conn.Close()
		return
	}

	strMsg := string(msgBytes)
	common.LogfDebug(strMsg)

	req := msg.WSJoinRequest{}
	if err := json.Unmarshal(msgBytes, &req); err != nil {
		common.LogfDebug("failed to unmarshal message: %s", err.Error())
		conn.Close()
		return
	}

	s.createOrJoinRoom(conn, &req)
}

func (s *WSService) createOrJoinRoom(conn *websocket.Conn, req *msg.WSJoinRequest) {
	if ok, messages := s.validate(req); !ok {
		WriteResponse(conn, msg.NewRejectJoinResponse(messages...))
		conn.Close()
		return
	}

	room, ok := s.Rooms[req.RoomCode]
	if !ok && req.AllowCreate {
		if len(s.Rooms) >= s.config.RoomCapacity {
			WriteResponse(conn, msg.NewRejectJoinResponse("failed to create room"))
			conn.Close()
			return
		}

		newRoom := NewWSRoom(req, s.roomCloseChan)
		newRoom.Start()

		s.Rooms[req.RoomCode] = newRoom
		room = newRoom
	} else {
		if !ok {
			WriteResponse(conn, msg.NewRejectJoinResponse("room not found"))
			conn.Close()
			return
		}

		if room.Locked || len(room.Clients) >= room.Capacity {
			WriteResponse(conn, msg.NewRejectJoinResponse("join not allowed"))
			conn.Close()
			return
		}
	}

	room.AddUser(conn, req.PlayerName)
}

func (s *WSService) validate(req *msg.WSJoinRequest) (ok bool, messages []string) {
	messages = make([]string, 0)

	if req == nil {
		messages = append(messages, "no request provided")
		return len(messages) == 0, messages
	}

	if ok, err := regexp.Match("^[A-Z]{4,10}$", []byte(req.RoomCode)); !ok || err != nil {
		messages = append(messages, "invalid room code format.  must be [A-Z], 4-10 characters long.")
	}

	if req.PlayerName == "" {
		messages = append(messages, "no player name set")
	}

	if req.AllowCreate {
		if req.GameID == "" {
			messages = append(messages, "no game id")
		}

		if req.RoomSize == 0 {
			messages = append(messages, "no room size")
		}
	}

	return len(messages) == 0, messages

}

func (s *WSService) hasSubProtocolHeader(r *http.Request) bool {
	for _, sub := range websocket.Subprotocols(r) {
		common.LogfDebug("Checking sub: %s\n", sub)
		if sub == s.config.Subprotocol {
			common.LogfDebug("Protocol match!")
			return true
		}
	}

	return false
}

func (s *WSService) Start() {
	go func() {
		for {
			room := <-s.roomCloseChan
			common.LogfDebug("removing room %s", room.Code)
			delete(s.Rooms, room.Code)
		}
	}()
}

func NewWSService(config ServerConfig) *WSService {
	return &WSService{
		config:        config,
		Rooms:         make(map[string]*WSRoom),
		roomCloseChan: make(chan *WSRoom),
	}
}
