package store

import (
	"errors"
	"reflect"

	"github.com/dchest/uniuri"
)

// Queues - Contain all the stored queues
var Queues []Queue

// Queue -
type Queue struct {
	ID            string
	Slug          string
	Group         Group
	RecurrentJobs []RecurrentJob
	Jobs          []Job
	Workers       []User
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
//                   If create flag set to true, the queue will be created
func FindQueueBySlug(slugToFind string, group Group, create bool) (Queue, error) {
	var queue Queue
	for _, value := range Queues {
		if value.Group.Slug == group.Slug && value.Slug == slugToFind {
			value.Group = group
			queue = value
		}
	}
	if create && queue.ID == "" {
		// if this queue don't exists, we create it
		queue.ID = uniuri.New()
		queue.Slug = slugToFind
		queue.Group = group
		Queues = append(Queues, queue)
	} else if queue.ID == "" {
		return queue, errors.New("unknown-queue")
	}
	return queue, nil
}

// Update -
func (queue Queue) Update() error {
	return queue.UpdateAndKeep([]string{})
}

// UpdateAndKeep - Will update a queue
func (queue Queue) UpdateAndKeep(toKeepFields []string) error {
	for key, value := range Queues {
		if value.ID == queue.ID {
			for _, field := range toKeepFields {
				reflect.ValueOf(&queue).Elem().FieldByName(field).Set(
					reflect.ValueOf(&value).Elem().FieldByName(field),
				)
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
