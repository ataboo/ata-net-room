package ws

import (
	"encoding/base64"
	"encoding/json"
	"time"
)

type RequestType string

const (
	JoinReq    RequestType = "join"
	LeaveReq   RequestType = "leave"
	LockReq    RequestType = "lock"
	UnlockReq  RequestType = "unlock"
	GameEvtReq RequestType = "game"
)

type ResponseType string

const (
	JoinReject    ResponseType = "join_reject"
	YouJoinRes    ResponseType = "join_you"
	PlayerJoinRes ResponseType = "join_player"
	LeaveRes      ResponseType = ResponseType(LeaveReq)
	LockRes       ResponseType = ResponseType(LockReq)
	UnlockRes     ResponseType = ResponseType(UnlockReq)
	GameEvtRes    ResponseType = ResponseType(GameEvtReq)
)

type ServerConfig struct {
	Host         string
	RoomCapacity int
	Subprotocol  string
}

type WSRequest struct {
	Type     RequestType `json:"type"`
	SendTime time.Time   `json:"send"`
	Payload  string      `json:"payload,omitempty"`
}

type WSJoinRequest struct {
	RoomCode    string `json:"room_code"`
	AllowCreate bool   `json:"create"`
	PlayerName  string `json:"player_name"`
	GameID      string `json:"game_id"`
	RoomSize    int    `json:"room_size"`
}

type WSResponse struct {
	Type      ResponseType `json:"type"`
	SendTime  time.Time    `json:"send"`
	RelayTime time.Time    `json:"relay,omitempty"`
	Payload   string       `json:"payload,omitempty"`
}

func (r *WSResponse) SetPayload(v interface{}) error {
	m, err := json.Marshal(v)
	if err != nil {
		return err
	}

	r.Payload = base64.RawStdEncoding.EncodeToString(m)
	return nil
}

func (r *WSRequest) SetPayload(v interface{}) error {
	m, err := json.Marshal(v)
	if err != nil {
		return err
	}

	r.Payload = base64.RawStdEncoding.EncodeToString(m)
	return nil
}
