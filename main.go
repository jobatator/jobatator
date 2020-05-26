package main

import (
	"github.com/lefuturiste/jobatator/pkg/server"
	"github.com/lefuturiste/jobatator/pkg/store"
	"github.com/lefuturiste/jobatator/pkg/utils"
)

func main() {
	utils.StartUptimeTimer()
	store.LoadConfigFromFile("./config.yml")
	go server.StartHTTPServer()
	server.Start()
}
