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
	err := queue.Update()
	if err != nil {
		ReturnError(cmd, err.Error())
		return
	}

	// dispatch old jobs
	go DispatchUniversal()

	ReturnString(cmd, "OK")
}

// Unsubscribe - If the client is a worker, he will use this cmd to unsubscribe from a queue
func Unsubscribe(cmd CmdInterface) {
	queue, err := store.FindQueueBySlug(cmd.Parts[1], cmd.User.CurrentGroup, false)
	if err != nil {
		ReturnError(cmd, err.Error())
		return
	}
	var workers []store.User
	for _, user := range queue.Workers {
		if user.Addr != cmd.User.Addr {
			workers = append(workers, user)
		}
	}
	queue.Workers = workers
	err = queue.Update()
	if err != nil {
		ReturnError(cmd, err.Error())
		return
	}
	ReturnString(cmd, "OK")
}
