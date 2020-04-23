package commands

import (
	"encoding/json"
	"time"

	"github.com/lefuturiste/jobatator/pkg/utils"
)

// Dispatch -
/*
This function will itterate on all the queue
for each queue:
	for each 'pending' job:
		for each workers in the queue:
			if worker is availaible:
				send job

	for each 'in-progress' job:
		if the job has expired:
			for each workers in the queue:
				if worker is availaible:
					send job

	for each 'errored' job:
		for each workers in the queue:
			if worker is availaible:
				send job
*/

// DispatchData -
type DispatchData struct {
	Job   utils.Job
	Debug string
}

// DispatchUniversal -
func DispatchUniversal() {
	for _, queue := range utils.Queues {
		for _, job := range queue.Jobs {

			// check if a job has expired, if the job expired, set as pending
			if job.Type == utils.JobInProgress && job.Attempts < 3 {
				duration := time.Since(job.StartedProcessingAt)
				if duration.Minutes() > 5 {
					// we consider the job as expired if the job started processing 5 min ago
					job.State = utils.JobPending
				}
			}

			if job.State == utils.JobPending || job.State == utils.JobErrored {
				// if the job is pending or errored we send the job to a available worker
				for _, worker := range queue.Workers {
					if worker.Status == utils.WorkerAvailable {
						// send the job
						dispatchData := DispatchData{
							Job:   job,
							Debug: "Dispatch data",
						}
						data, _ := json.Marshal(dispatchData)
						worker.Conn.Write(data)
						worker.Conn.Write([]byte("\n"))
					}
				}
			}
		}
	}
}

// Dispatch -
func Dispatch(cmd utils.CmdInterface) {
	DispatchUniversal()
}
