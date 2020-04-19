package server

import (
	"encoding/json"

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
	Job utils.Job
}

// Dispatch -
func Dispatch() {
	for _, queue := range utils.Queues {
		for _, job := range queue.Jobs {
			if job.Type == utils.JobPending {
				for _, worker := range queue.Workers {
					if worker.Status == utils.WorkerAvailable {
						// send the job
						dispatchData := DispatchData{
							Job: job,
						}
						data, _ := json.Marshal(dispatchData)
						worker.Conn.Write(data)
						worker.Conn.Write([]byte("\r\n"))
					}
				}
			}
			// check if a job has expired, if the job expired, set as pending
			// if job.Type == utils.JobInProgress {
			// 	duration := time.Since(job.StartedProcessingAt)
			// 	if duration.Minutes() > 5 {
			// 		// we consider the job as expired if the job started processing 5 min ago
			// 		job.State = utils.JobPending
			// 	}
			// }
		}
	}
}
