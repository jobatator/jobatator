package commands

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/jobatator/jobatator/pkg/store"
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

// StopServer - Will exit the server process, only work if TEST_MODE env var or config key 'test_mode' is set to true
//              Warning: this feature is really dangerous and open serious security risks. Use it with cares.
func StopServer(cmd CmdInterface) {
	if store.Options.TestMode {
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
		if value == group.Slug || value == "*" {
			isAllowed = true
		}
	}
	if !isAllowed {
		return group, errors.New("forbidden-group")
	}
	user.CurrentGroup = group
	store.UpdateSession(user)
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
	// we update the session
	cmd.User.Username = user.Username
	cmd.User.Password = user.Password
	cmd.User.Groups = user.Groups
	store.UpdateSession(cmd.User)
	ReturnString(cmd, "Welcome!")
}

// SessionOutput - Represent the current session
type SessionOutput struct {
	Username     string
	Addr         string
	CurrentGroup string
	IsWorker     bool
	Groups       []string
}

// CurrentSession - Get current session meta data
func CurrentSession(cmd CmdInterface) {
	user := store.FindSession(cmd.Conn)
	jsonStr, _ := json.Marshal(SessionOutput{
		Username:     user.Username,
		Addr:         user.Addr,
		CurrentGroup: user.CurrentGroup.Slug,
		IsWorker:     user.Status != "",
		Groups:       user.Groups,
	})
	ReturnString(cmd, string(jsonStr))
}
