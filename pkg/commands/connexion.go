package commands

import (
	"strconv"

	"github.com/lefuturiste/jobatator/pkg/utils"
)

// Ping -
func Ping(cmd utils.CmdInterface) {
	utils.ReturnString(cmd, "PONG")
}

// Quit -
func Quit(cmd utils.CmdInterface) {
	cmd.Conn.Close()
}

// UseGroup - will switch the session on a specific group
func UseGroup(cmd utils.CmdInterface) {
	if len(cmd.Parts) != 2 {
		utils.ReturnError(cmd, "invalid-input")
		return
	}
	var group utils.Group
	for _, value := range utils.Options.Groups {
		if value.Slug == cmd.Parts[1] {
			group = value
		}
	}
	if group.Slug == "" {
		utils.ReturnError(cmd, "unknown-group")
		return
	}
	user := utils.FindSession(cmd)
	var isAllowed bool = false
	for _, value := range user.Groups {
		if value == group.Slug {
			isAllowed = true
		}
	}
	if !isAllowed {
		utils.ReturnError(cmd, "forbidden-group")
		return
	}
	for key, value := range utils.Sessions {
		if value.Addr == user.Addr {
			value.CurrentGroup = group
			utils.Sessions[key] = value
		}
	}
	utils.ReturnString(cmd, "OK")
}

// Auth -
func Auth(cmd utils.CmdInterface) {
	if len(cmd.Parts) != 3 {
		utils.ReturnError(cmd, "invalid-input")
		return
	}
	var user utils.User
	for _, val := range utils.Options.Users {
		if val.Username == cmd.Parts[1] && val.Password == cmd.Parts[2] {
			user = val
		}
	}
	if user.Username == "" {
		utils.ReturnError(cmd, "invalid-creds")
		return
	}
	utils.ReturnString(cmd, "Welcome!")
	if len(utils.Sessions) == 0 {
		utils.Sessions = make([]utils.User, 0)
	}
	// we add the user to the list of the sessions
	user.Addr = cmd.Conn.RemoteAddr().String()
	user.Conn = cmd.Conn
	utils.Sessions = append(utils.Sessions, user)
}

// Debug -
func Debug(cmd utils.CmdInterface) {
	utils.ReturnString(cmd, "== GROUPS ==")
	for _, group := range utils.Options.Groups {
		utils.ReturnString(cmd, "- slug: "+group.Slug)
	}

	utils.ReturnString(cmd, "\n== SESSION USERS ==")
	for _, user := range utils.Sessions {
		utils.ReturnString(cmd, "- username: "+user.Username)
		utils.ReturnString(cmd, "  currentGroup: "+user.CurrentGroup.Slug)
		utils.ReturnString(cmd, "  addr: "+user.Addr)
	}

	utils.ReturnString(cmd, "\n== QUEUES ==")
	for _, queue := range utils.Queues {
		utils.ReturnString(cmd, "- slug: "+queue.Slug)
		utils.ReturnString(cmd, "  jobs: "+strconv.FormatInt(int64(len(queue.Jobs)), 10))
		for _, job := range queue.Jobs {
			utils.ReturnString(cmd, "    - id: "+job.ID)
			utils.ReturnString(cmd, "      type: "+job.Type)
			utils.ReturnString(cmd, "      state: "+job.State)
			utils.ReturnString(cmd, "      payload: "+job.Payload)
		}
		utils.ReturnString(cmd, "  workers: "+strconv.FormatInt(int64(len(queue.Workers)), 10))
		for _, worker := range queue.Workers {
			utils.ReturnString(cmd, "    - addr: "+worker.Addr)
			utils.ReturnString(cmd, "      username: "+worker.Username)
		}

	}
	if len(utils.Queues) == 0 {
		utils.ReturnString(cmd, "No queues")
	}
}
