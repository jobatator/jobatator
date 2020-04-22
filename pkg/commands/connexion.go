package commands

import (
	"errors"
	"os"

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

// StopServer - Will exit the server process, only work if the env var STOP_POLICY is PUBLIC.
//              Warning: this feature is really dangerous and open serious secuity risks. Use it with cares.
func StopServer(cmd utils.CmdInterface) {
	if os.Getenv("STOP_POLICY") == "PUBLIC" {
		utils.NewLine(cmd)
		os.Exit(0)
	} else {
		utils.ReturnError(cmd, "forbidden")
	}
}

// UseGroupUniversal -
func UseGroupUniversal(groupRaw string, user utils.User) (utils.Group, error) {
	var group utils.Group
	for _, value := range utils.Options.Groups {
		if value.Slug == groupRaw {
			group = value
		}
	}
	if group.Slug == "" {
		return group, errors.New("unknown-group")
	}
	var isAllowed bool = false
	for _, value := range user.Groups {
		if value == group.Slug {
			isAllowed = true
		}
	}
	if !isAllowed {
		return group, errors.New("forbidden-group")
	}
	user.CurrentGroup = group
	utils.UpdateUser(user)
	return group, nil
}

// UseGroup - will switch the session on a specific group
func UseGroup(cmd utils.CmdInterface) {
	if len(cmd.Parts) != 2 {
		utils.ReturnError(cmd, "invalid-input")
		return
	}
	user := utils.FindSession(cmd)
	_, err := UseGroupUniversal(cmd.Parts[1], user)

	if err == nil {
		utils.ReturnString(cmd, "OK")
	} else {
		utils.ReturnError(cmd, err.Error())
	}
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
