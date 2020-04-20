package main

import (
	"github.com/lefuturiste/jobatator/pkg/server"
	"github.com/lefuturiste/jobatator/pkg/utils"
)

func main() {
	utils.LoadConfigFromFile("./config.yml")
	go server.StartHTTPServer()
	server.Start()
}
