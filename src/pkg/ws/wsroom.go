package ws

type WSRoom struct {
	Code     string
	Capacity int
	Locked   bool
	Client   []WSClient
}
