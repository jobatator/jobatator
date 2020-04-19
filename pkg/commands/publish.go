package commands

import (
	"github.com/dchest/uniuri"
	"github.com/lefuturiste/jobatator/pkg/utils"
)

// Publish - Will add a job on a queue PUBLISH queue_name job_type payload
func Publish(cmd utils.CmdInterface) {
	if len(cmd.Parts) != 4 {
		utils.ReturnError(cmd, "invalid-input")
		return
	}
	user := utils.FindSession(cmd)
	if user.CurrentGroup.Slug == "" {
		utils.ReturnError(cmd, "group-non-selected")
		return
	}
	if len(utils.Queues) == 0 {
		utils.Queues = make([]utils.Queue, 0)
	}
	var queue utils.Queue
	// find the queue
	for _, value := range utils.Queues {
		if value.Slug == cmd.Parts[1] {
			queue = value
		}
	}
	if queue.Slug == "" {
		// if this queue don't exists, we create it
		queue.Slug = cmd.Parts[1]
		queue.Group = user.CurrentGroup
		utils.Queues = append(utils.Queues, queue)
	}
	if len(queue.Jobs) == 0 {
		queue.Jobs = make([]utils.Job, 0)
	}
	var job utils.Job
	job.ID = uniuri.New()
	job.State = utils.JobPending
	job.Type = cmd.Parts[2]
	job.Payload = cmd.Parts[3]
	queue.Jobs = append(queue.Jobs, job)

	// update the queue state into the db
	for key, value := range utils.Queues {
		if value.Slug == queue.Slug {
			utils.Queues[key] = queue
		}
	}
	// if a worker is availaible, notify a worker

	go Dispatch()

	utils.ReturnString(cmd, "OK")
}
