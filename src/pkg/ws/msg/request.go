package msg

type RequestType int

const (
	GameEvtReq RequestType = 0
	JoinReq    RequestType = 1
	LeaveReq   RequestType = 2
	LockReq    RequestType = 3
	UnlockReq  RequestType = 4
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
