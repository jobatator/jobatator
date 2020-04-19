package server

import (
	"fmt"
	"net"
	"strings"
)

var state map[string]string
var password string

func handleClient(conn net.Conn) {
	fmt.Println("New client:", conn.RemoteAddr().String())
	var message bool
	var input string
	var componentIndex int
	var components map[int]string
	if state == nil {
		state = make(map[string]string)
	}
	for {
		message = true
		buf := make([]byte, 1024)
		_, err := conn.Read(buf)
		if err != nil {
			conn.Close()
			break
		}
		for index, value := range buf {
			if 0 == value {
				// 10 mean linefeed, so end of the line
				input += string(buf)[0:index]

				if strings.Count(input, "\r\n") > 1 {
					componentIndex = 0
					components = make(map[int]string)
					for _, component := range strings.Split(input, "\r\n") {
						if len(component) == 0 {
							break
						}
						if !(component[0:1] == "*" || component[0:1] == "$") {
							components[componentIndex] = component
							componentIndex++
						}
					}
				} else {
					components = parseCommand(input[0 : len(input)-1])
				}
				if len(components) > 0 {
					fmt.Println("New cmd: ", components)
					//var name string = strings.ToUpper(components[0])

					//fmt.Println(components)
				}

				input = ""
				message = false
				break
			}
		}
		if message {
			input = input + string(buf)
		}
	}
}
