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
	go serverLoop()
}

// Start - Start the TCP server
func Start() {
	serverLoop()
}

func serverLoop() {
	var host string = utils.Options.Host
	var port string = strconv.FormatInt(int64(utils.Options.Port), 10)
	l, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer l.Close()
	fmt.Println("Listening on " + host + ":" + port)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		go handleClient(conn)
	}
}
