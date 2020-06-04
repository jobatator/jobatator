package server

import (
	"net"
	"os"
	"strconv"

	"github.com/jobatator/jobatator/pkg/store"
	log "github.com/sirupsen/logrus"
)

// StartAsync - Start the TCP server with a go routine
func StartAsync() {
	listener := prepareServer()
	go serverLoop(listener)
}

// Start - Start the TCP server
func Start() {
	listener := prepareServer()
	serverLoop(listener)
}

func prepareServer() net.Listener {
	var host string = store.Options.Host
	var port string = strconv.FormatInt(int64(store.Options.Port), 10)
	listener, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		log.Error("Error listening: ", err.Error())
		os.Exit(1)
	}
	log.Info("Listening on " + host + ":" + port)
	return listener
}

func serverLoop(listener net.Listener) {
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Error("Error accepting connexion: ", err.Error())
			os.Exit(1)
		}
		go handleClient(conn)
	}
}
