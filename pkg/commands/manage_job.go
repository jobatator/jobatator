package commands

import (
	"encoding/json"
	"time"

	"github.com/jobatator/jobatator/pkg/store"
)

// ListJobs - List all jobs in a queue
func ListJobs(cmd CmdInterface) {
	queue, err := store.FindQueueBySlug(cmd.Parts[1], cmd.User.CurrentGroup, false)
	if err != nil {
		ReturnError(cmd, err.Error())
		return
	}

	rawJSON, _ := json.Marshal(queue.Jobs)
	ReturnString(cmd, string(rawJSON))
}

// UpdateJob - Update the state of the job
// Arguments: JOB_ID, JOB_STATUS
func UpdateJob(cmd CmdInterface) {
	job, err := store.FindJob(cmd.Parts[1])
	if err != nil {
		ReturnError(cmd, err.Error())
		return
	}

	if cmd.Parts[2] == store.JobInProgress {
		job.State = store.JobInProgress
		job.StartedProcessingAt = time.Now()
	} else if cmd.Parts[2] == store.JobDone {
		job.State = store.JobDone
		job.EndProcessingAt = time.Now()
	} else if cmd.Parts[2] == store.JobErrored {
		job.State = store.JobErrored
		job.Attempts = job.Attempts + 1
	} else {
		ReturnError(cmd, "unknown-state")
		return
	}

	store.UpdateSession(cmd.User)
	job.Update()

	delay := -1
	if cmd.Parts[2] == store.JobInProgress {
		// set the worker as busy
		cmd.User.Status = store.WorkerBusy
	} else {
		// free this worker
		cmd.User.Status = store.WorkerAvailable
		// see if this worker can work again
		delay = 2 // in seconds
		if job.State == store.JobErrored {
			delay = 1800 // wait 30 minutes before trying again
		}
		if store.Options.DelayPolicy == "IGNORE" {
			delay = 0
		}
	}

	// update the worker status
	for key, worker := range job.Queue.Workers {
		if worker.Addr == cmd.User.Addr {
			job.Queue.Workers[key] = cmd.User
		}
	}
	job.Queue.UpdateAndKeep([]string{"Jobs"})

	if delay != -1 {
		go DispatchUniversalWithDelay(delay)
	}

	if job.State == store.JobDone {
		go job.Expire(store.Options.JobTimeout)
	}

	ReturnString(cmd, "OK")
}

// DeleteJob - Delete a job
func DeleteJob(cmd CmdInterface) {
	job, err := store.FindJob(cmd.Parts[1])
	if err != nil {
		ReturnError(cmd, err.Error())
		return
	}
	err = job.Delete()
	if err != nil {
		ReturnError(cmd, err.Error())
		return
	}

	ReturnString(cmd, "OK")
}
