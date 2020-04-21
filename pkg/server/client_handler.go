package server

import (
	"fmt"
	"net"
	"strings"

	"github.com/lefuturiste/jobatator/pkg/utils"
)

func handleClient(conn net.Conn) {
	fmt.Println("New client:", conn.RemoteAddr().String())
	var message bool
	var input string
	var componentIndex int
	var components map[int]string
	var user utils.User
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
					var name string = strings.ToUpper(components[0])

					for _, val := range utils.Sessions {
						if val.Addr == conn.RemoteAddr().String() {
							user = val
						}
					}
					if user.Username == "" && name != "AUTH" {
						conn.Write([]byte("Err: not-logged"))
						conn.Write([]byte("\r\n"))
					} else {
						cmd := utils.CmdInterface{
							Parts: components,
							Conn:  conn,
						}
						if CmdMap[name] == nil {
							conn.Write([]byte("Err: unkown-command"))
							conn.Write([]byte("\r\n"))
						} else {
							CmdMap[name].(func(utils.CmdInterface))(cmd)
						}
					}
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
	if user.Username != "" {
		currentAddr := conn.RemoteAddr().String()
		// remove user from the session array
		var newSessions []utils.User
		for _, val := range utils.Sessions {
			if currentAddr != val.Addr {
				newSessions = append(newSessions, val)
			}
		}
		utils.Sessions = newSessions

		// remove user from all the worker inside the queues
		for key, queue := range utils.Queues {
			var newWorkers []utils.User
			for _, worker := range queue.Workers {
				if currentAddr != worker.Addr {
					newWorkers = append(newWorkers, worker)
				}
			}
			utils.Queues[key].Workers = newWorkers
		}
	}
}
