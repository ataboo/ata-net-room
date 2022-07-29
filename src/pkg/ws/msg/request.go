package msg

type RequestType int

const (
	GameEvtReq RequestType = 0
	LeaveReq   RequestType = 1
	LockReq    RequestType = 2
	UnlockReq  RequestType = 3
)

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
