package utils

import "time"

// Sessions - all the sessions
var Sessions []User

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
}

// Queue -
type Queue struct {
	Group   Group
	Slug    string
	Jobs    []Job
	Workers []User
}

// WorkerBusy -
const WorkerBusy = "busy"

// WorkerAvailable -
const WorkerAvailable = "available"

// Queues -
var Queues []Queue

// FindSession - will return a user object for this session
func FindSession(cmd CmdInterface) User {
	var user User
	for _, value := range Sessions {
		if value.Addr == cmd.Conn.RemoteAddr().String() {
			user = value
		}
	}
	return user
}

// UpdateUser -
func UpdateUser(user User) {
	for key, value := range Sessions {
		if value.Addr == user.Addr {
			Sessions[key] = user
		}
	}
}
