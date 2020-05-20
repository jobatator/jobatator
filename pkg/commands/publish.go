package commands

import (
	"github.com/dchest/uniuri"
	"github.com/lefuturiste/jobatator/pkg/store"
)

// PublishUniversal - Will add a job on a queue PUBLISH queue_name job_type payload
func PublishUniversal(parts map[int]string, user store.User) (string, error) {
	queue, _ := store.FindQueueBySlug(parts[1], user.CurrentGroup, true)

	// create a job
	var job store.Job
	job.ID = uniuri.New()
	job.State = store.JobPending
	job.Type = parts[2]
	job.Payload = parts[3]
	queue.Jobs = append(queue.Jobs, job)
	queue.Update(false)

	// dispatch the created job
	go DispatchUniversal()

	return job.ID, nil
}

// Publish - Cli interface
func Publish(cmd CmdInterface) {
	result, err := PublishUniversal(cmd.Parts, cmd.User)

	if err != nil {
		ReturnError(cmd, err.Error())
	} else {
		ReturnString(cmd, "OK#"+result)
	}
}
