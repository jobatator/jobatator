package server

import (
	"net"
	"strings"

	cmds "github.com/jobatator/jobatator/pkg/commands"
	"github.com/jobatator/jobatator/pkg/store"
	log "github.com/sirupsen/logrus"
)

func handleClient(conn net.Conn) {
	currentAddr := conn.RemoteAddr().String()
	log.Debug(currentAddr + " - new client")
	var message bool
	var input string
	var componentIndex int
	components := make(map[int]string)
	cmd := cmds.CmdInterface{
		Conn: conn,
	}
	// create a session
	var user store.User
	user.Addr = currentAddr
	user.Conn = conn
	store.Sessions = append(store.Sessions, user)

	for {
		message = true
		buf := make([]byte, 1024)
		_, err := conn.Read(buf)
		if err != nil {
			//log.Error("Failed to read from socket client:", err)
			conn.Close()
			break
		}
		for index, value := range buf {
			if 0 == value {
				// 10 mean linefeed, so end of the line
				input += string(buf)[0:index]
				if strings.Count(input, "\r\n") > 1 {
					componentIndex = 0
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
					log.Debug(currentAddr+" - cmd - ", components)
					var name string = strings.ToUpper(components[0])
					var cmdDefinition cmds.CmdDefinition
					cmdDefinition.Args = -1
					cmdDefinition.UseGroup = false
					cmdDefinition.RequireAuth = false
					cmd.Parts = components
					for _, value := range cmds.CmdMap {
						if value.Name == name {
							cmdDefinition = value
						}
					}
					if cmdDefinition.Name == "" && name == "HELP" {
						cmdDefinition.Name = "HELP"
						cmdDefinition.Handler = cmds.Help
					}
					if cmdDefinition.Name == "" {
						cmds.ReturnError(cmd, "unknown-command")
					} else {
						handleCommand(cmdDefinition, cmd)
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

	// user disconnected, we need to delete the user session and remove as a worker in queues
	// remove user from the session array
	var newSessions []store.User
	for _, val := range store.Sessions {
		if currentAddr != val.Addr {
			newSessions = append(newSessions, val)
		}
	}
	store.Sessions = newSessions

	// remove user from all the workers inside the queues
	for key, queue := range store.Queues {
		var newWorkers []store.User
		for _, worker := range queue.Workers {
			if currentAddr != worker.Addr {
				newWorkers = append(newWorkers, worker)
			}
		}
		store.Queues[key].Workers = newWorkers
	}
	log.Debug(currentAddr + " - quit")
}

func handleCommand(definition cmds.CmdDefinition, cmd cmds.CmdInterface) {
	user := store.FindSession(cmd.Conn)
	if definition.RequireAuth && user.Username == "" {
		cmds.ReturnError(cmd, "not-logged")
		return
	}
	if definition.UseGroup {
		if user.CurrentGroup.Slug == "" {
			cmds.ReturnError(cmd, "group-non-selected")
			return
		}
	}
	cmd.User = user
	if definition.Args != -1 && len(cmd.Parts) != definition.Args+1 {
		cmds.ReturnError(cmd, "wrong-numbers-arguments")
		return
	}
	definition.Handler(cmd)
}
