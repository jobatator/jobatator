package commands

import (
	"time"

	"github.com/lefuturiste/jobatator/pkg/utils"
)

// UpdateJob - Update the state of the job
// ACCEPT QUEUE_NAME JOB_ID JOB_STATUS
func UpdateJob(cmd utils.CmdInterface) {
	if len(cmd.Parts) != 4 {
		utils.ReturnError(cmd, "invalid-input")
		return
	}
	user := utils.FindSession(cmd)
	if user.CurrentGroup.Slug == "" {
		utils.ReturnError(cmd, "group-non-selected")
		return
	}
	var queue utils.Queue
	// find the queue
	for _, value := range utils.Queues {
		if value.Slug == cmd.Parts[1] {
			queue = value
		}
	}
	if queue.Slug == "" {
		utils.ReturnError(cmd, "unknown-queue")
		return
	}
	var job utils.Job
	var jobKey int
	for key, value := range queue.Jobs {
		if value.ID == cmd.Parts[2] {
			jobKey = key
			job = value
		}
	}
	if job.ID == "" {
		utils.ReturnError(cmd, "unknown-job")
		return
	}

	if cmd.Parts[3] == utils.JobInProgress {
		job.State = utils.JobInProgress
		job.StartedProcessingAt = time.Now()
	} else if cmd.Parts[3] == utils.JobDone {
		job.State = utils.JobDone
		job.EndProcessingAt = time.Now()
	} else if cmd.Parts[3] == utils.JobErrored {
		job.State = utils.JobErrored
		job.Attempts = job.Attempts + 1
	} else {
		utils.ReturnError(cmd, "unknown-state")
		return
	}

	if cmd.Parts[3] == utils.JobInProgress {
		// set the worker as busy
		user.Status = utils.WorkerBusy
	} else {
		// free this worker
		user.Status = utils.WorkerAvailable
		// see if this worker can work again
		delay := 2 // in seconds
		if job.State == utils.JobErrored {
			delay = 1800 // wait 30 minutes before trying again
		}
		if utils.Options.DelayPolicy == "IGNORE" {
			delay = 0
		}
		go DispatchUniversalWithDelay(delay)
	}
	utils.UpdateUser(user)

	queue.Jobs[jobKey] = job

	for key, val := range queue.Workers {
		if val.Addr == user.Addr {
			queue.Workers[key] = user
		}
	}

	// update the job state into the db
	for key, value := range utils.Queues {
		if value.Slug == queue.Slug {
			utils.Queues[key] = queue
		}
	}

	utils.ReturnString(cmd, "OK")
}
