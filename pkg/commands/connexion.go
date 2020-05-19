package commands

import (
	"errors"
	"os"

	"github.com/lefuturiste/jobatator/pkg/store"
)

// Ping -
func Ping(cmd CmdInterface) {
	ReturnString(cmd, "PONG")
}

// Quit -
func Quit(cmd CmdInterface) {
	NewLine(cmd)
	cmd.Conn.Close()
}

// StopServer - Will exit the server process, only work if the env var STOP_POLICY is PUBLIC.
//              Warning: this feature is really dangerous and open serious secuity risks. Use it with cares.
func StopServer(cmd CmdInterface) {
	if os.Getenv("STOP_POLICY") == "PUBLIC" {
		NewLine(cmd)
		os.Exit(0)
	} else {
		ReturnError(cmd, "forbidden")
	}
}

// UseGroupUniversal -
func UseGroupUniversal(groupRaw string, user store.User) (store.Group, error) {
	var group store.Group
	for _, value := range store.Options.Groups {
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
	store.UpdateUser(user)
	return group, nil
}

// UseGroup - will switch the session on a specific group
func UseGroup(cmd CmdInterface) {
	_, err := UseGroupUniversal(cmd.Parts[1], cmd.User)

	if err == nil {
		ReturnString(cmd, "OK")
	} else {
		ReturnError(cmd, err.Error())
	}
}

// Auth -
func Auth(cmd CmdInterface) {
	var user store.User
	for _, val := range store.Options.Users {
		if val.Username == cmd.Parts[1] && val.Password == cmd.Parts[2] {
			user = val
		}
	}
	if user.Username == "" {
		ReturnError(cmd, "invalid-creds")
		return
	}
	ReturnString(cmd, "Welcome!")
	if len(store.Sessions) == 0 {
		store.Sessions = make([]store.User, 0)
	}
	// we add the user to the list of the sessions
	user.Addr = cmd.Conn.RemoteAddr().String()
	user.Conn = cmd.Conn
	store.Sessions = append(store.Sessions, user)
}
