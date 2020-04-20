package commands

import (
	"encoding/json"
	"strconv"

	"github.com/lefuturiste/jobatator/pkg/utils"
)

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
			utils.ReturnString(cmd, "      status: "+worker.Status)
		}

	}
	if len(utils.Queues) == 0 {
		utils.ReturnString(cmd, "No queues")
	}
}

// DebugOutput -
type DebugOutput struct {
	Queues   []utils.Queue
	Groups   []utils.Group
	Sessions []utils.User
	Users    []utils.User
	Host     string
	Port     int
}

// DebugJSON -
func DebugJSON(cmd utils.CmdInterface) {
	var debubOutput DebugOutput
	debubOutput.Queues = utils.Queues
	debubOutput.Groups = utils.Options.Groups
	debubOutput.Sessions = utils.Sessions
	debubOutput.Users = utils.Options.Users
	debubOutput.Host = utils.Options.Host
	debubOutput.Port = utils.Options.Port
	rawJSON, _ := json.Marshal(debubOutput)
	utils.ReturnString(cmd, string(rawJSON))
}
