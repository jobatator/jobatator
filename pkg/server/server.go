package server

import (
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/lefuturiste/jobatator/pkg/utils"
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
	var host string = utils.Options.Host
	var port string = strconv.FormatInt(int64(utils.Options.Port), 10)
	listener, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	fmt.Println("Listening on " + host + ":" + port)
	return listener
}

func serverLoop(listener net.Listener) {
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		go handleClient(conn)
	}
}
