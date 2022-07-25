package ws

type JoinData struct {
	RoomCode    string `json:"room_code"`
	AllowCreate bool   `json:"create"`
	PlayerName  string `json:"player_name"`
	GameID      string `json:"game_id"`
	RoomSize    int    `json:"room_size"`
}

type ServerConfig struct {
	Host         string
	RoomCapacity int
	Subprotocol  string
}
