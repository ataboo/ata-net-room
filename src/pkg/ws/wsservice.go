package ws

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/ataboo/ata-net-room/pkg/common"
	"github.com/gorilla/websocket"
)

const (
	MaxMessageSize = 512
	ReadWait       = 3 * time.Second
	WriteWait      = 3 * time.Second
)

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
	config ServerConfig
	Rooms  map[string]WSRoom
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
	mType, msg, err := conn.ReadMessage()
	if mType != websocket.BinaryMessage || err != nil {
		common.LogfDebug("expected binary message")
		conn.Close()
	}

	req := WSJoinRequest{}
	if err := json.Unmarshal(msg, &req); err != nil {
		common.LogfDebug("failed to unmarshal message: %s", err.Error())
		conn.Close()
	}

	s.createOrJoinRoom(conn, &req)
}

func (s *WSService) createOrJoinRoom(conn *websocket.Conn, req *WSJoinRequest) {
	if ok, messages := s.validate(req); !ok {
		res := WSResponse{
			Type:      JoinReject,
			SendTime:  time.Now().UnixMilli(),
			RelayTime: time.Now().UnixMilli(),
			Payload:   MustEncodeBase64Payload(messages),
			ID:        GenUniqueID(),
		}

		conn.SetWriteDeadline(time.Now().Add(WriteWait))
		conn.WriteJSON(res)

		conn.Close()
		return
	}

	if req.AllowCreate {
		newRoom := WSRoom{}
	}
}

func (s *WSService) validate(req *WSJoinRequest) (ok bool, messages []string) {
	messages = make([]string, 0)

	if req == nil {
		messages = append(messages, "no request provided")
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
		fmt.Printf("Checking sub: %s\n", sub)
		if sub == s.config.Subprotocol {
			return true
		}
	}

	return false
}

func NewWSService(config ServerConfig) *WSService {
	return &WSService{
		config: config,
		Rooms:  make(map[string]WSRoom),
	}
}
