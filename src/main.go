package main

import (
	"github.com/ataboo/ata-net-room/pkg/rest"
	"github.com/ataboo/ata-net-room/pkg/ws"
)

func main() {
	server := rest.NewRestServer(ws.ServerConfig{
		Host:         ":3000",
		RoomCapacity: 10,
		Subprotocol:  "atanet_v1",
	})

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
