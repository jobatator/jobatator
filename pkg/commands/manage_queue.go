package commands

import (
	"encoding/json"

	"github.com/lefuturiste/jobatator/pkg/store"
)

// ListQueues - List all the queues in a group
func ListQueues(cmd CmdInterface) {
	var queues []store.Queue
	// filter all the queues inside a particular group
	for _, queue := range store.Queues {
		if queue.Group.Slug == cmd.User.CurrentGroup.Slug {
			queues = append(queues, queue)
		}
	}
	rawJSON, _ := json.Marshal(queues)
	ReturnString(cmd, string(rawJSON))
}

// DeleteQueue - Delete a queue
func DeleteQueue(cmd CmdInterface) {
	queue, err := store.FindQueueBySlug(cmd.Parts[1], cmd.User.CurrentGroup)
	if err != nil {
		ReturnError(cmd, err.Error())
		return
	}
	err = queue.Delete()
	if err != nil {
		ReturnError(cmd, err.Error())
		return
	}
	ReturnString(cmd, "OK")
}
