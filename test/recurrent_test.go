package test

import (
	"bufio"
	"encoding/json"
	"testing"

	"github.com/lefuturiste/jobatator/pkg/commands"
	"github.com/lefuturiste/jobatator/pkg/store"
	"github.com/magiconair/properties/assert"
)

var recurrentWorkerReady bool = false

func recurrentWorkerClient(t *testing.T) {
	conn := getConn()
	buf := bufio.NewReader(conn)
	doAuthStuff(conn)

	send(conn, "SUBSCRIBE recurrent")
	assert.Equal(t, readReply(buf), "OK")

	recurrentWorkerReady = true

	reply := readReply(buf)
	var dispatchData commands.DispatchData
	json.Unmarshal([]byte(reply), &dispatchData)
	assert.Equal(t, dispatchData.Job.Type, "my_recurrent_type")
	assert.Equal(t, dispatchData.Job.Payload, "")

	send(conn, "UPDATE_JOB "+dispatchData.Job.ID+" done")
	reply = readReply(buf)
	assert.Equal(t, reply, "OK")

	recurrentWorkerReady = true
}

func TestRecurrent(t *testing.T) {
	startServer()
	go recurrentWorkerClient(t)
	conn := getConn()
	buf := bufio.NewReader(conn)
	doAuthStuff(conn)

	// wait for the worker client to initialize
	for !recurrentWorkerReady {
	}
	recurrentWorkerReady = false

	// try to declare with a wrong cron expression
	send(conn, "RECURRENT_JOB recurrent my_recurrent_type 'invalid'")
	reply := readReply(buf)
	assert.Equal(t, reply[0:18], "Err: invalid-cron;")

	// declare a recurrent job
	send(conn, "RECURRENT_JOB recurrent my_recurrent_type '* * * * *'")
	reply = readReply(buf)
	// we get the id
	assert.Equal(t, reply, "OK#1")

	// assert that the recurrent job was created
	debug := getDebug(conn, buf)
	var queue store.Queue
	for _, value := range debug.Queues {
		if value.Slug == "recurrent" {
			queue = value
		}
	}
	assert.Equal(t, queue.Slug, "recurrent")
	assert.Equal(t, queue.RecurrentJobs[0].Type, "my_recurrent_type")
	assert.Equal(t, queue.RecurrentJobs[0].CronExpression, "* * * * *")
	assert.Equal(t, queue.RecurrentJobs[0].EntryID, 1)

	// assert that a job has been immidialty created
	assert.Equal(t, queue.Slug, "recurrent")
	assert.Equal(t, queue.Jobs[0].Type, "my_recurrent_type")
	assert.Equal(t, queue.Jobs[0].Payload, "")

	// wait for the worker client to process the job
	for !recurrentWorkerReady {
	}
	recurrentWorkerReady = false
}
