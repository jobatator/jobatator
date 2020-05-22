package store

import (
	"errors"
	"time"
)

// JobPending - the job is in the queue ready to be taken by one of the available worker
const JobPending = "pending"

// JobInProgress - the job was taken by a worker and being worked out
const JobInProgress = "in-progress"

// JobErrored - the job was processed but issues happen so it need a rerun
const JobErrored = "errored"

// JobDone - the job was processed without issues by the worker
const JobDone = "done"

// Job -
type Job struct {
	ID                  string
	State               string
	Type                string
	Payload             string
	Attempts            int
	StartedProcessingAt time.Time
	EndProcessingAt     time.Time
	Queue               Queue
	Group               Group
}

// FindJob - find a job in all the queues
func FindJob(id string) (Job, error) {
	var job Job
	// find the job
	for _, queue := range Queues {
		for _, job := range queue.Jobs {
			if job.ID == id {
				queue.Jobs = make([]Job, 0)
				job.Queue = queue
				return job, nil
			}
		}
	}
	return job, errors.New("unknown-job")
}

// Update - Update a job
func (job Job) Update() error {
	queue, err := FindQueue(job.Queue.ID)
	if err != nil {
		return err
	}
	var jobKey int = -1
	for key, value := range queue.Jobs {
		if value.ID == job.ID {
			jobKey = key
		}
	}
	if jobKey == -1 {
		return errors.New("unknown-job")
	}
	job.Queue = Queue{}
	queue.Jobs[jobKey] = job
	queue.Update()
	return nil
}

// Delete - Will delete a job
func (job Job) Delete() error {
	queue, err := FindQueue(job.Queue.ID)
	if err != nil {
		return err
	}
	var newJobs []Job
	for _, value := range queue.Jobs {
		if value.ID != job.ID {
			newJobs = append(newJobs, value)
		}
	}
	queue.Jobs = newJobs
	err = queue.Update()
	if err != nil {
		return err
	}
	return nil
}

// Expire - Will delete a job after a timeout in seconds
func (job Job) Expire(timeout int) error {
	time.Sleep(time.Duration(timeout) * time.Second)
	err := job.Delete()
	if err != nil {
		return err
	}
	return nil
}
