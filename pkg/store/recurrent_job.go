package store

import "errors"

// RecurrentJob -
type RecurrentJob struct {
	EntryID        int
	Queue          Queue
	Type           string
	CronExpression string
}

// FindRecurrentJob - Find a reccurent job by id
func FindRecurrentJob(id int) (RecurrentJob, error) {
	var job RecurrentJob
	for _, queue := range Queues {
		for _, job := range queue.RecurrentJobs {
			if job.EntryID == id {
				queue.Jobs = make([]Job, 0)
				queue.RecurrentJobs = make([]RecurrentJob, 0)
				job.Queue = queue
				return job, nil
			}
		}
	}
	return job, errors.New("unknown-recurrent-job")
}

// FindRecurrentJobByType - Find a reccurent job by id
func FindRecurrentJobByType(queue Queue, jobType string) (RecurrentJob, error) {
	var defaultJob RecurrentJob
	for _, job := range queue.RecurrentJobs {
		if job.Type == jobType {
			return job, nil
		}
	}
	return defaultJob, errors.New("unknown-recurrent-job")
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
