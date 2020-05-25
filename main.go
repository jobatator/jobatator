package main

import (
	"github.com/lefuturiste/jobatator/pkg/server"
	"github.com/lefuturiste/jobatator/pkg/store"
)

func main() {
	store.StartUptimeTimer()
	store.LoadConfigFromFile("./config.yml")
	go server.StartHTTPServer()
	server.Start()
}
