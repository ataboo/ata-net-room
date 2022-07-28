package ws

import (
	"encoding/base64"
	"encoding/json"

	"github.com/google/uuid"
)

type RequestType int

const (
	GameEvtReq RequestType = 0
	JoinReq    RequestType = 1
	LeaveReq   RequestType = 2
	LockReq    RequestType = 3
	UnlockReq  RequestType = 4
)

type ResponseType int

const (
	GameEvtRes    ResponseType = 0
	JoinReject    ResponseType = 1
	YouJoinRes    ResponseType = 2
	PlayerJoinRes ResponseType = 3
	LeaveRes      ResponseType = 4
	LockRes       ResponseType = 5
	UnlockRes     ResponseType = 6
)

type ServerConfig struct {
	Host         string
	RoomCapacity int
	Subprotocol  string
}

type WSRequest struct {
	Type     RequestType `json:"type"`
	SendTime int64       `json:"send"`
	ID       string      `json:"id,omitempty"`
	Name     string      `json:"name,omitempty"`
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
	SendTime  int64        `json:"send"`
	RelayTime int64        `json:"relay"`
	ID        string       `json:"id,omitempty"`
	Name      string       `json:"name,omitempty"`
	Payload   string       `json:"payload,omitempty"`
}

func MustEncodeBase64Payload(v interface{}) string {
	encoded, err := EncodeBase64Payload(v)
	if err != nil {
		panic(err)
	}

	return encoded
}

func EncodeBase64Payload(v interface{}) (payload string, err error) {
	m, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	return base64.RawStdEncoding.EncodeToString(m), nil
}

func GenUniqueID() string {
	return uuid.NewString()
}
