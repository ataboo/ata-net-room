package msg

import (
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/google/uuid"
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

type PlayerIDPayload struct {
	PlayerID        int   `json:"id"`
	OtherPlayersIDs []int `json:"players"`
}

type WSResponse struct {
	Type      ResponseType `json:"type"`
	SendTime  int64        `json:"send"`
	RelayTime int64        `json:"relay"`
	ID        string       `json:"id,omitempty"`
	Name      string       `json:"name,omitempty"`
	Payload   string       `json:"payload,omitempty"`
}

func NewRejectJoinResponse(messages ...string) WSResponse {
	return WSResponse{
		Type:      JoinReject,
		SendTime:  time.Now().UnixMilli(),
		RelayTime: time.Now().UnixMilli(),
		Payload:   MustEncodeBase64Payload(messages),
		ID:        GenUniqueID(),
	}
}

func NewJoinResponse(youJoin bool, payload PlayerIDPayload) WSResponse {
	resType := PlayerJoinRes
	if youJoin {
		resType = YouJoinRes
	}

	return WSResponse{
		Type:      resType,
		SendTime:  time.Now().UnixMilli(),
		RelayTime: time.Now().UnixMilli(),
		Payload:   MustEncodeBase64Payload(payload),
		ID:        GenUniqueID(),
	}
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
