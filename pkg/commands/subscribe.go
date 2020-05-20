package commands

import (
	"github.com/lefuturiste/jobatator/pkg/store"
)

// Subscribe - If the client is a worker, he will use this cmd to subscribe to a queue
func Subscribe(cmd CmdInterface) {
	queue, _ := store.FindQueueBySlug(cmd.Parts[1], cmd.User.CurrentGroup, true)

	// we register the user as worker in this queue
	cmd.User.Status = store.WorkerAvailable
	queue.Workers = append(queue.Workers, cmd.User)
	queue.Update(true)

	// dispatch old jobs
	go DispatchUniversal()

	ReturnString(cmd, "OK")
}
