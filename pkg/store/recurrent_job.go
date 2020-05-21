package store

import "errors"

// RecurrentJob -
type RecurrentJob struct {
	EntryID        int
	Queue          Queue
	Type           string
	CronExpression string
}

// ListRecurrentJob - List all reccurent jobs in a queue
func ListRecurrentJob(queue Queue) []RecurrentJob {
	var jobs []RecurrentJob
	return jobs
}

// FindRecurrentJob - Find a reccurent job by id
func FindRecurrentJob() RecurrentJob {
	var job RecurrentJob
	return job
}

// FindRecurrentJobByType - Find a reccurent job by id
func FindRecurrentJobByType(queue Queue, jobType string) (RecurrentJob, error) {
	var defaultJob RecurrentJob
	for _, job := range queue.RecurrentJobs {
		if job.Type == jobType {
			return job, nil
		}
	}
	return defaultJob, errors.New("unknown-job")
}

// Update - Update a reccurent job
func (job RecurrentJob) Update() error {
	queue, err := FindQueue(job.Queue.ID)
	if err != nil {
		return err
	}
	var jobKey int = -1
	for key, value := range queue.RecurrentJobs {
		if value.EntryID == job.EntryID {
			jobKey = key
		}
	}
	if jobKey == -1 {
		return errors.New("unknown-recurrent-job")
	}
	job.Queue = Queue{}
	queue.RecurrentJobs[jobKey] = job
	queue.Update()
	return nil
}

// Delete - Delete a reccurent job
func (job RecurrentJob) Delete() error {
	return nil
}
