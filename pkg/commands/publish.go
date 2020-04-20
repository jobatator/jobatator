package commands

import (
	"errors"

	"github.com/dchest/uniuri"
	"github.com/lefuturiste/jobatator/pkg/utils"
)

// PublishUniversal - Will add a job on a queue PUBLISH queue_name job_type payload
func PublishUniversal(parts map[int]string, user utils.User) (bool, error) {
	if len(parts) != 4 {
		return false, errors.New("invalid-input")
	}
	if user.CurrentGroup.Slug == "" {
		return false, errors.New("group-non-selected")
	}
	if len(utils.Queues) == 0 {
		utils.Queues = make([]utils.Queue, 0)
	}
	var queue utils.Queue
	// find the queue
	for _, value := range utils.Queues {
		if value.Slug == parts[1] {
			queue = value
		}
	}
	if queue.Slug == "" {
		// if this queue don't exists, we create it
		queue.Slug = parts[1]
		queue.Group = user.CurrentGroup
		utils.Queues = append(utils.Queues, queue)
	}
	if len(queue.Jobs) == 0 {
		queue.Jobs = make([]utils.Job, 0)
	}
	var job utils.Job
	job.ID = uniuri.New()
	job.State = utils.JobPending
	job.Type = parts[2]
	job.Payload = parts[3]
	queue.Jobs = append(queue.Jobs, job)

	// update the queue state into the db
	for key, value := range utils.Queues {
		if value.Slug == queue.Slug {
			utils.Queues[key] = queue
		}
	}
	// if a worker is availaible, notify a worker

	go DispatchUniversal()

	return true, nil
}

// Publish - Cli interface
func Publish(cmd utils.CmdInterface) {
	user := utils.FindSession(cmd)
	result, err := PublishUniversal(cmd.Parts, user)

	if result {
		utils.ReturnString(cmd, "OK")
	} else {
		utils.ReturnError(cmd, err.Error())
	}
}
