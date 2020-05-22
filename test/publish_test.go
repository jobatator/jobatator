package test

import (
	"bufio"
	"encoding/json"
	"net"
	"strconv"
	"testing"

	"github.com/dchest/uniuri"
	"github.com/lefuturiste/jobatator/pkg/commands"
	"github.com/lefuturiste/jobatator/pkg/server"
	"github.com/lefuturiste/jobatator/pkg/store"
	"github.com/magiconair/properties/assert"
)

const testConfig = `
port: 8963
host: "0.0.0.0"
groups:
    - slug: "group1"
    - slug: "group2"
    - slug: "group3"

users:
    - username: "user1"
      password: "pass1"
      groups: ["group1"]
    - username: "user2"
      password: "pass1"
      groups: ["group2", "group3"]

delay_policy: "IGNORE"
test_mode: true
log_level: "DEBUG"
`

type FooObject struct {
	ID        string
	Name      string
	Slug      string
	Something string
}

type JobArgs struct {
	UserID           string
	FileURL          string
	IsSomething      bool
	CountOfSomething int
	List             []FooObject
}

var fakeJobArgs JobArgs

func getFakeJobArgs() JobArgs {
	var myJobArgs JobArgs
	myJobArgs.UserID = "ID"
	myJobArgs.FileURL = "https://example.com"
	myJobArgs.IsSomething = true
	myJobArgs.CountOfSomething = 9942
	myJobArgs.List = make([]FooObject, 0)

	for i := 1; i < 10; i++ {
		var fooObject FooObject
		fooObject.ID = uniuri.New()
		fooObject.Name = "Object " + strconv.FormatInt(int64(i), 10)
		fooObject.Slug = "object-" + strconv.FormatInt(int64(i), 10)
		fooObject.Something = "lel-" + uniuri.New() + "-lel"
		myJobArgs.List = append(myJobArgs.List, fooObject)
	}
	return myJobArgs
}

func getDebug(conn net.Conn, buf *bufio.Reader) commands.DebugOutput {
	send(conn, "DEBUG_JSON")
	reply := readReply(buf)
	var debugData commands.DebugOutput
	json.Unmarshal([]byte(reply), &debugData)
	return debugData
}

var secondClientReady bool = false

func secondClient(t *testing.T) {
	conn := getConn()
	buf := bufio.NewReader(conn)

	send(conn, "AUTH user1 pass1")
	reply := readReply(buf)
	assert.Equal(t, reply, "Welcome!")

	send(conn, "USE_GROUP group1")
	reply = readReply(buf)
	assert.Equal(t, reply, "OK")

	send(conn, "SUBSCRIBE default")
	reply = readReply(buf)
	assert.Equal(t, reply, "OK")
	debug := getDebug(conn, buf)
	assert.Equal(t, len(debug.Queues), 1)
	assert.Equal(t, debug.Queues[0].Slug, "default")

	// the second client is ready
	secondClientReady = true

	reply = readReply(buf)
	var dispatchData commands.DispatchData
	json.Unmarshal([]byte(reply), &dispatchData)
	assert.Equal(t, dispatchData.Job.Type, "job.type")
	jsonRaw, _ := json.Marshal(fakeJobArgs)
	assert.Equal(t, dispatchData.Job.Payload, string(jsonRaw))

	send(conn, "UPDATE_JOB "+dispatchData.Job.ID+" errored")
	reply = readReply(buf)
	assert.Equal(t, reply, "OK")

	// waiting to received one again the job because it is errored
	reply = readReply(buf)
	json.Unmarshal([]byte(reply), &dispatchData)
	assert.Equal(t, dispatchData.Job.State, "errored")
	debug = getDebug(conn, buf)
	assert.Equal(t, debug.Queues[0].Jobs[0].State, store.JobErrored)
	assert.Equal(t, debug.Queues[0].Workers[0].Status, store.WorkerAvailable)

	send(conn, "UPDATE_JOB "+dispatchData.Job.ID+" in-progress")
	reply = readReply(buf)
	assert.Equal(t, reply, "OK")
	debug = getDebug(conn, buf)
	assert.Equal(t, debug.Queues[0].Jobs[0].State, store.JobInProgress)
	assert.Equal(t, debug.Queues[0].Workers[0].Status, store.WorkerBusy)

	send(conn, "UPDATE_JOB "+dispatchData.Job.ID+" done")
	reply = readReply(buf)
	assert.Equal(t, reply, "OK")
	debug = getDebug(conn, buf)
	assert.Equal(t, debug.Queues[0].Jobs[0].State, store.JobDone)
	assert.Equal(t, debug.Queues[0].Workers[0].Status, store.WorkerAvailable)

	secondClientReady = true
}

func TestJob(t *testing.T) {
	store.LoadConfigFromString(testConfig)
	server.StartAsync()

	go secondClient(t)

	conn := getConn()
	buf := bufio.NewReader(conn)

	// try to authenticate with bad credidentials
	send(conn, "AUTH user1 pass2")
	reply := readReply(buf)
	assert.Equal(t, reply[0:3], "Err")

	// try to authenticate with good credidentials
	send(conn, "AUTH user1 pass1")
	reply = readReply(buf)
	assert.Equal(t, reply, "Welcome!")

	// check if the connexion is alive
	send(conn, "PING")
	reply = readReply(buf)
	assert.Equal(t, reply, "PONG")

	// try to use a forbidden group
	send(conn, "USE_GROUP group2")
	reply = readReply(buf)
	assert.Equal(t, reply[0:3], "Err")

	// try to use a allowed group
	send(conn, "USE_GROUP group1")
	reply = readReply(buf)
	assert.Equal(t, reply, "OK")

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
	reply = readReply(buf)
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
