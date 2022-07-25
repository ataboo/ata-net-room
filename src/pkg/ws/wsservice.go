package ws

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
		log.Printf("WS Error: %d, %s\n", status, reason.Error())
	},
	Subprotocols: []string{"atanet_v1"},
}

type WSService struct {
	config ServerConfig
	Rooms  map[string]WSRoom
}

func (s *WSService) HandleJoin(w http.ResponseWriter, r *http.Request) {
	// if r.Method != http.MethodPost {
	// 	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	// 	return
	// }

	// data := JoinData{}
	// if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
	// 	http.Error(w, "malformed request", http.StatusBadRequest)
	// 	return
	// }

	if !s.hasSubProtocolHeader(r) {
		http.Error(w, "invalid protocol", http.StatusBadRequest)
		return
	}

	h := http.Header{}
	h.Set("Sec-Websocket-Protocol", s.config.Subprotocol)

	conn, err := upgrader.Upgrade(w, r, h)
	if err != nil {
		http.Error(w, "failed to upgrade", http.StatusUnprocessableEntity)
		return
	}

	conn.WriteMessage(websocket.BinaryMessage, []byte("Hello WS!"))

	conn.Close()

	// log.Printf("got join: %+v", data)

	// TODO room stuff, add connection

	// http.Error(w, "join denied", http.StatusUnauthorized)
	// return
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
