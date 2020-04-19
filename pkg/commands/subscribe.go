package commands

import (
	"github.com/lefuturiste/jobatator/pkg/utils"
)

// Subscribe - If the client is a worker, he will use this cmd to subscribe to a queue
func Subscribe(cmd utils.CmdInterface) {
	if len(cmd.Parts) != 2 {
		utils.ReturnError(cmd, "invalid-input")
		return
	}
	user := utils.FindSession(cmd)
	if user.CurrentGroup.Slug == "" {
		utils.ReturnError(cmd, "group-non-selected")
		return
	}
	// find the queue
	var queue utils.Queue
	var queueKey int
	// find the queue
	for key, value := range utils.Queues {
		if value.Slug == cmd.Parts[1] {
			queue = value
			queueKey = key
		}
	}
	if queue.Slug == "" {
		// if this queue don't exists, we create it
		queue.Slug = cmd.Parts[1]
		queue.Group = user.CurrentGroup
		utils.Queues = append(utils.Queues, queue)
	}
	if len(queue.Workers) == 0 {
		queue.Workers = make([]utils.User, 0)
	}
	user.Status = utils.WorkerAvailable
	// we register the user as worker in this queue
	queue.Workers = append(queue.Workers, user)
	utils.Queues[queueKey] = queue

	go Dispatch()

	utils.ReturnString(cmd, "OK")
}
