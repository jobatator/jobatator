package commands

import (
	"errors"

	"github.com/dchest/uniuri"
	"github.com/lefuturiste/jobatator/pkg/store"
)

// PublishUniversal - Will add a job on a queue PUBLISH queue_name job_type payload
func PublishUniversal(parts map[int]string, user store.User) (bool, error) {
	if len(parts) != 4 {
		return false, errors.New("invalid-input")
	}
	if user.CurrentGroup.Slug == "" {
		return false, errors.New("group-non-selected")
	}
	if len(store.Queues) == 0 {
		store.Queues = make([]store.Queue, 0)
	}
	var queue store.Queue
	// find the queue
	for _, value := range store.Queues {
		if value.Slug == parts[1] {
			queue = value
		}
	}
	if queue.Slug == "" {
		// if this queue don't exists, we create it
		queue.ID = uniuri.New()
		queue.Slug = parts[1]
		queue.Group = user.CurrentGroup
		store.Queues = append(store.Queues, queue)
	}
	if len(queue.Jobs) == 0 {
		queue.Jobs = make([]store.Job, 0)
	}
	var job store.Job
	job.ID = uniuri.New()
	job.State = store.JobPending
	job.Type = parts[2]
	job.Payload = parts[3]
	queue.Jobs = append(queue.Jobs, job)

	// update the queue state into the db
	for key, value := range store.Queues {
		if value.Slug == queue.Slug {
			store.Queues[key] = queue
		}
	}
	// if a worker is availaible, notify a worker

	go DispatchUniversal()

	return true, nil
}

// Publish - Cli interface
func Publish(cmd CmdInterface) {
	result, err := PublishUniversal(cmd.Parts, cmd.User)

	if result {
		ReturnString(cmd, "OK")
	} else {
		ReturnError(cmd, err.Error())
	}
}
