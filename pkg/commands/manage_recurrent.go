package commands

import (
	"encoding/json"
	"strconv"

	"github.com/jobatator/jobatator/pkg/store"
)

// ListRecurrentJobs - List all recurrent jobs in a queue
func ListRecurrentJobs(cmd CmdInterface) {
	queue, err := store.FindQueueBySlug(cmd.Parts[1], cmd.User.CurrentGroup, false)
	if err != nil {
		ReturnError(cmd, err.Error())
		return
	}

	rawJSON, _ := json.Marshal(queue.RecurrentJobs)
	ReturnString(cmd, string(rawJSON))
}

// DeleteRecurrentJobs - Delete a recurrent job with a entry id in a queue
func DeleteRecurrentJobs(cmd CmdInterface) {
	entryID, _ := strconv.Atoi(cmd.Parts[1])
	job, err := store.FindRecurrentJob(entryID)
	if err != nil {
		ReturnError(cmd, err.Error())
		return
	}
	err = job.Delete()
	if err != nil {
		ReturnError(cmd, err.Error())
		return
	}

	ReturnString(cmd, "OK")
}
