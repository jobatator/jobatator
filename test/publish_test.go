package test

import (
	"bufio"
	"encoding/json"
	"testing"

	"github.com/jobatator/jobatator/pkg/commands"
	"github.com/jobatator/jobatator/pkg/store"
	"github.com/magiconair/properties/assert"
)

var secondClientReady bool = false

func secondClient(t *testing.T) {
	conn := getConn()
	doAuthStuff(conn)
	buf := bufio.NewReader(conn)

	// try to subscribe to the default queue
	send(conn, "SUBSCRIBE default")
	reply := readReply(buf)
	assert.Equal(t, reply, "OK")

	// assert that the queue is there
	debug := getDebug(conn, buf)
	assert.Equal(t, len(debug.Queues), 1)
	assert.Equal(t, debug.Queues[0].Slug, "default")

	// assert that the worker is there
	assert.Equal(t, len(debug.Queues[0].Workers), 1)
	assert.Equal(t, debug.Queues[0].Workers[0].Username, "user1")
	assert.Equal(t, debug.Queues[0].Workers[0].CurrentGroup.Slug, "group1")

	// the second client is ready
	secondClientReady = true

	// read the job dispatch data
	reply = readReply(buf)
	var dispatchData commands.DispatchData
	json.Unmarshal([]byte(reply), &dispatchData)
	assert.Equal(t, dispatchData.Job.Type, "job.type")

	// assert the payload is as expected
	jsonRaw, _ := json.Marshal(fakeJobArgs)
	assert.Equal(t, dispatchData.Job.Payload, string(jsonRaw))

	// try to set the job as updated
	send(conn, "UPDATE_JOB "+dispatchData.Job.ID+" errored")
	reply = readReply(buf)
	assert.Equal(t, reply, "OK")

	// waiting to received one again the job because it is errored
	reply = readReply(buf)
	json.Unmarshal([]byte(reply), &dispatchData)
	assert.Equal(t, dispatchData.Job.State, "errored")

	// assert tthat the job is set to a errored state and thus the worker to a available state
	debug = getDebug(conn, buf)
	assert.Equal(t, debug.Queues[0].Jobs[0].State, store.JobErrored)
	assert.Equal(t, debug.Queues[0].Workers[0].Status, store.WorkerAvailable)

	// try to update the job to an 'in-progress' state
	send(conn, "UPDATE_JOB "+dispatchData.Job.ID+" in-progress")
	reply = readReply(buf)
	assert.Equal(t, reply, "OK")

	// assert that the job is set as in progress and thus the worker as busy
	debug = getDebug(conn, buf)
	assert.Equal(t, debug.Queues[0].Jobs[0].State, store.JobInProgress)
	assert.Equal(t, debug.Queues[0].Workers[0].Status, store.WorkerBusy)

	// try to set the job as done
	send(conn, "UPDATE_JOB "+dispatchData.Job.ID+" done")
	reply = readReply(buf)
	assert.Equal(t, reply, "OK")

	// assert that the job is updated
	debug = getDebug(conn, buf)
	assert.Equal(t, debug.Queues[0].Jobs[0].State, store.JobDone)
	assert.Equal(t, debug.Queues[0].Workers[0].Status, store.WorkerAvailable)

	// try to unsubscribe from the queue
	send(conn, "UNSUBSCRIBE default")
	reply = readReply(buf)
	assert.Equal(t, reply, "OK")

	// assert that the worker is deleted
	debug = getDebug(conn, buf)
	assert.Equal(t, len(debug.Queues[0].Workers), 0)

	secondClientReady = true
}

func TestJob(t *testing.T) {
	startServer()

	go secondClient(t)

	conn := getConn()
	doAuthStuff(conn)
	buf := bufio.NewReader(conn)

	// wait for the second client to initialize (auth & use group)
	for !secondClientReady {
	}
	secondClientReady = false

	// create dummy job payload
	fakeJobArgs = getFakeJobArgs()
	jsonRaw, _ := json.Marshal(fakeJobArgs)
	jsonStr := string(jsonRaw)

	// publish the job
	send(conn, "PUBLISH default job.type '"+jsonStr+"'")
	reply := readReply(buf)
	// to get the id use reply[3:]
	assert.Equal(t, reply[0:3], "OK#")

	// assert that the job was created
	debug := getDebug(conn, buf)
	assert.Equal(t, debug.Queues[0].Slug, "default")
	assert.Equal(t, debug.Queues[0].Jobs[0].Type, "job.type")
	assert.Equal(t, debug.Queues[0].Jobs[0].Payload, jsonStr)

	// wait for the second client to process the job
	for !secondClientReady {
	}
	secondClientReady = false

	// check if the connexion is still alive
	send(conn, "PING")
	reply = readReply(buf)
	assert.Equal(t, reply, "PONG")

	// try list the queues
	send(conn, "LIST_QUEUES")
	reply = readReply(buf)
	var queues []store.Queue
	json.Unmarshal([]byte(reply), &queues)
	assert.Equal(t, len(queues), 1)

	// try list the jobs
	send(conn, "LIST_JOBS "+debug.Queues[0].Slug)
	reply = readReply(buf)
	var jobs []store.Job
	json.Unmarshal([]byte(reply), &jobs)
	assert.Equal(t, len(jobs), 1)
	// assert that the job is done
	assert.Equal(t, jobs[0].State, store.JobDone)
	assert.Equal(t, jobs[0].Type, "job.type")

	// try to delete a job
	send(conn, "DELETE_JOB "+debug.Queues[0].Jobs[0].ID)
	beforeCount := len(debug.Queues[0].Jobs)
	reply = readReply(buf)
	assert.Equal(t, reply, "OK")

	// assert that the job has been deleted
	debug = getDebug(conn, buf)
	assert.Equal(t, len(debug.Queues[0].Jobs), beforeCount-1)

	// try to delete a queue
	send(conn, "DELETE_QUEUE "+debug.Queues[0].Slug)
	reply = readReply(buf)
	assert.Equal(t, reply, "OK")

	// assert that the queue has been deleted
	debug = getDebug(conn, buf)
	assert.Equal(t, 0, len(debug.Queues))
}
