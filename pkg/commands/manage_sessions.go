package commands

import (
	"encoding/json"

	"github.com/jobatator/jobatator/pkg/store"
)

// ListSessions - List all sessions in the current group
func ListSessions(cmd CmdInterface) {
	sessions := make([]store.User, 0)
	for _, session := range store.Sessions {
		if session.CurrentGroup.Slug == cmd.User.CurrentGroup.Slug {
			sessions = append(sessions, session)
		}
	}

	rawJSON, _ := json.Marshal(sessions)
	ReturnString(cmd, string(rawJSON))
}
