package store

import (
	"errors"
)

// Queues - Contain all the stored queues
var Queues []Queue

// Queue -
type Queue struct {
	ID      string
	Slug    string
	Group   Group
	Jobs    []Job
	Workers []User
}

// ListQueue - Return all queues in a group
func ListQueue(group Group) []Queue {
	var queues []Queue
	for _, queue := range Queues {
		if queue.Group.Slug == group.Slug {
			queue.Group = group
			queues = append(queues, queue)
		}
	}
	return queues
}

// FindQueue - Will find a queue with a id
func FindQueue(id string) (Queue, error) {
	for _, queue := range Queues {
		if queue.ID == id {
			return queue, nil
		}
	}
	return Queue{}, errors.New("unknown-queue")
}

// FindQueueBySlug - Will find a queue with a slug in a group
func FindQueueBySlug(slugToFind string, group Group) (Queue, error) {
	for _, queue := range Queues {
		if queue.Group.Slug == group.Slug && queue.Slug == slugToFind {
			queue.Group = group
			return queue, nil
		}
	}
	return Queue{}, errors.New("unknown-queue")
}

// Update - Will update a queue
func (queue Queue) Update(keepJobs bool) error {
	for key, value := range Queues {
		if value.ID == queue.ID {
			if keepJobs {
				queue.Jobs = value.Jobs
			}
			Queues[key] = queue
			return nil
		}
	}
	return errors.New("unknown-queue")
}

// Delete - Will delete a queue
func (queue Queue) Delete() error {
	var newQueues []Queue
	for _, value := range Queues {
		if value.ID != queue.ID {
			newQueues = append(newQueues, value)
		}
	}
	Queues = newQueues
	return nil
}
