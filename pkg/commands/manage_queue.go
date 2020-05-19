package commands

import "github.com/lefuturiste/jobatator/pkg/store"

// ListQueues - List all the queues in a group
func ListQueues(cmd CmdInterface) {
	if cmd.User.CurrentGroup.Slug == "" {
		ReturnError(cmd, "group-non-selected")
		return
	}
	var queues []store.Queue
	// filter all the queues inside a particular group
	for _, queue := range store.Queues {
		if queue.Group.Slug == cmd.User.CurrentGroup.Slug {
			queues = append(queues, queue)
		}
	}

	ReturnString(cmd, "OK")
}
