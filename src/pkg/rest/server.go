package rest

import (
	"encoding/json"
	"net/http"

	"github.com/ataboo/ata-net-room/pkg/ws"
)

func NewRestServer(config ws.ServerConfig) *RestServer {
	mux := http.NewServeMux()
	wsService := ws.NewWSService(config)
	server := &RestServer{
		config:    config,
		mux:       mux,
		wsService: wsService,
	}

	server.initHandlers()

	return server
}

func (s *RestServer) initHandlers() {
	s.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		WriteJSON(w, r, map[string]string{"status": "OK"}, http.StatusOK)
	})

	s.mux.HandleFunc("/join", func(w http.ResponseWriter, r *http.Request) {
		s.wsService.HandleJoin(w, r)
	})

	s.mux.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		serverInfo := ServerInfo{
			RoomCount:    len(s.wsService.Rooms),
			RoomCapacity: s.config.RoomCapacity,
		}

		WriteJSON(w, r, serverInfo, http.StatusOK)
	})
}

type RestServer struct {
	config    ws.ServerConfig
	mux       *http.ServeMux
	wsService *ws.WSService
}

func (s *RestServer) ListenAndServe() error {
	return http.ListenAndServeTLS(s.config.Host, "server.crt", "server.key", s.mux)
	//return http.ListenAndServe(s.config.Host, s.mux)
}

func WriteJSON(w http.ResponseWriter, r *http.Request, data interface{}, status int) {
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
}
