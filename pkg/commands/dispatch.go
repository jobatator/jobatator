package commands

import (
	"encoding/json"
	"fmt"

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
	fmt.Println("New dispatch batch")
	for _, queue := range utils.Queues {
		fmt.Println(queue)
		for _, job := range queue.Jobs {
			fmt.Println(job)
			if job.State == utils.JobPending {
				fmt.Println("- One pendind job")
				for _, worker := range queue.Workers {
					fmt.Println(worker)
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
