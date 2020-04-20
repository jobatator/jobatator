package commands

import (
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
	user.CurrentGroup = group
	utils.UpdateUser(user)
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
