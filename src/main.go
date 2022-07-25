package main

import "github.com/ataboo/ata-net-room/pkg/http"

func main() {
	foo := &http.MyHttp{}

	if foo == nil {
		panic("nil but why?")
	}
}
