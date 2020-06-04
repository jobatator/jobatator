package main

import (
	"github.com/jobatator/jobatator/pkg/server"
	"github.com/jobatator/jobatator/pkg/store"
	"github.com/jobatator/jobatator/pkg/utils"
)

func main() {
	utils.StartUptimeTimer()
	store.LoadConfigFromFile("./config.yml")
	go server.StartHTTPServer()
	server.Start()
}
