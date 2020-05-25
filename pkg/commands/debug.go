package commands

import (
	"encoding/json"
	"strconv"

	"github.com/lefuturiste/jobatator/pkg/store"
)

// Debug -
func Debug(cmd CmdInterface) {
	if !store.Options.TestMode {
		ReturnError(cmd, "test-mode-disabled")
		return
	}
	ReturnString(cmd, "== GROUPS ==")
	for _, group := range store.Options.Groups {
		ReturnString(cmd, "- slug: "+group.Slug)
	}

	ReturnString(cmd, "\n== SESSION USERS ==")
	for _, user := range store.Sessions {
		ReturnString(cmd, "- username: "+user.Username)
		ReturnString(cmd, "  currentGroup: "+user.CurrentGroup.Slug)
		ReturnString(cmd, "  addr: "+user.Addr)
	}

	ReturnString(cmd, "\n== QUEUES ==")
	for _, queue := range store.Queues {
		ReturnString(cmd, "- slug: "+queue.Slug)
		ReturnString(cmd, "  id: "+queue.ID)
		ReturnString(cmd, "  jobs: "+strconv.FormatInt(int64(len(queue.Jobs)), 10))
		for _, job := range queue.Jobs {
			ReturnString(cmd, "    - id: "+job.ID)
			ReturnString(cmd, "      type: "+job.Type)
			ReturnString(cmd, "      state: "+job.State)
			ReturnString(cmd, "      payload: "+job.Payload)
		}
		ReturnString(cmd, "  recurrent_jobs: "+strconv.FormatInt(int64(len(queue.RecurrentJobs)), 10))
		for _, job := range queue.RecurrentJobs {
			ReturnString(cmd, "    - id: "+strconv.FormatInt(int64(job.EntryID), 10))
			ReturnString(cmd, "      type: "+job.Type)
			ReturnString(cmd, "      expression: "+job.CronExpression)
		}
		ReturnString(cmd, "  workers: "+strconv.FormatInt(int64(len(queue.Workers)), 10))
		for _, worker := range queue.Workers {
			ReturnString(cmd, "    - addr: "+worker.Addr)
			ReturnString(cmd, "      username: "+worker.Username)
			ReturnString(cmd, "      status: "+worker.Status)
		}

	}
	if len(store.Queues) == 0 {
		ReturnString(cmd, "No queues")
	}
}

// DebugOutput -
type DebugOutput struct {
	Queues   []store.Queue
	Groups   []store.Group
	Sessions []store.User
	Users    []store.User
	Host     string
	Port     int
}

// DebugJSON -
func DebugJSON(cmd CmdInterface) {
	if !store.Options.TestMode {
		ReturnError(cmd, "test-mode-disabled")
		return
	}
	var debubOutput DebugOutput
	debubOutput.Queues = store.Queues
	debubOutput.Groups = store.Options.Groups
	debubOutput.Sessions = store.Sessions
	debubOutput.Users = store.Options.Users
	debubOutput.Host = store.Options.Host
	debubOutput.Port = store.Options.Port
	rawJSON, _ := json.Marshal(debubOutput)
	ReturnString(cmd, string(rawJSON))
}

// DebugParts -
func DebugParts(cmd CmdInterface) {
	if !store.Options.TestMode {
		ReturnError(cmd, "test-mode-disabled")
		return
	}
	var parts []string
	for i := 0; i < len(cmd.Parts); i++ {
		parts = append(parts, cmd.Parts[i])
	}
	rawJSON, _ := json.Marshal(parts)
	ReturnString(cmd, string(rawJSON))
}
