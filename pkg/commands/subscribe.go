package commands

import (
	"github.com/dchest/uniuri"
	"github.com/lefuturiste/jobatator/pkg/store"
)

// Subscribe - If the client is a worker, he will use this cmd to subscribe to a queue
func Subscribe(cmd CmdInterface) {
	if cmd.User.CurrentGroup.Slug == "" {
		ReturnError(cmd, "group-non-selected")
		return
	}
	// find the queue
	var queue store.Queue
	var queueKey int
	// find the queue
	for key, value := range store.Queues {
		if value.Slug == cmd.Parts[1] {
			queue = value
			queueKey = key
		}
	}
	if queue.Slug == "" {
		// if this queue don't exists, we create it
		queue.ID = uniuri.New()
		queue.Slug = cmd.Parts[1]
		queue.Group = cmd.User.CurrentGroup
		store.Queues = append(store.Queues, queue)
	}
	if len(queue.Workers) == 0 {
		queue.Workers = make([]store.User, 0)
	}
	cmd.User.Status = store.WorkerAvailable
	// we register the user as worker in this queue
	queue.Workers = append(queue.Workers, cmd.User)
	store.Queues[queueKey] = queue

	go DispatchUniversal()

	ReturnString(cmd, "OK")
}
