package main

import (
	"bufio"
	"encoding/json"
	"net"
	"strconv"
	"strings"
	"testing"

	"github.com/dchest/uniuri"
	"github.com/lefuturiste/jobatator/pkg/commands"
	"github.com/lefuturiste/jobatator/pkg/server"
	"github.com/lefuturiste/jobatator/pkg/utils"
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
`

func getConn() net.Conn {
	tcpAddr, _ := net.ResolveTCPAddr("tcp", "localhost:8963")
	conn, _ := net.DialTCP("tcp", nil, tcpAddr)
	return conn
}

func readReply(buf *bufio.Reader) string {
	reply, _ := buf.ReadString('\n')
	return strings.Replace(reply, "\n", "", 1)
}

func send(conn net.Conn, str string) {
	conn.Write([]byte(str + "\n"))
}

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

	secondClientReady = true

	reply = readReply(buf)
	var dispatchData commands.DispatchData
	json.Unmarshal([]byte(reply), &dispatchData)
	assert.Equal(t, dispatchData.Job.Type, "job.type")

	jsonRaw, _ := json.Marshal(fakeJobArgs)
	assert.Equal(t, dispatchData.Job.Payload, string(jsonRaw))

	send(conn, "UPDATE_JOB default "+dispatchData.Job.ID+" errored")
	reply = readReply(buf)
	assert.Equal(t, reply, "OK")

	reply = readReply(buf)
	json.Unmarshal([]byte(reply), &dispatchData)
	assert.Equal(t, dispatchData.Job.State, "errored")
	debug = getDebug(conn, buf)
	assert.Equal(t, debug.Queues[0].Jobs[0].State, utils.JobErrored)
	assert.Equal(t, debug.Queues[0].Workers[0].Status, utils.WorkerAvailable)

	send(conn, "UPDATE_JOB default "+dispatchData.Job.ID+" in-progress")
	reply = readReply(buf)
	assert.Equal(t, reply, "OK")
	debug = getDebug(conn, buf)
	assert.Equal(t, debug.Queues[0].Jobs[0].State, utils.JobInProgress)
	assert.Equal(t, debug.Queues[0].Workers[0].Status, utils.WorkerBusy)

	send(conn, "UPDATE_JOB default "+dispatchData.Job.ID+" done")
	reply = readReply(buf)
	assert.Equal(t, reply, "OK")
	debug = getDebug(conn, buf)
	assert.Equal(t, debug.Queues[0].Jobs[0].State, utils.JobDone)
	assert.Equal(t, debug.Queues[0].Workers[0].Status, utils.WorkerAvailable)

	secondClientReady = true
}

func TestJob(t *testing.T) {
	utils.LoadConfigFromString(testConfig)
	server.StartAsync()

	go secondClient(t)

	conn := getConn()
	buf := bufio.NewReader(conn)

	send(conn, "AUTH user1 pass2")
	reply := readReply(buf)
	assert.Equal(t, reply[0:3], "Err")

	send(conn, "AUTH user1 pass1")
	reply = readReply(buf)
	assert.Equal(t, reply, "Welcome!")

	send(conn, "PING")
	reply = readReply(buf)
	assert.Equal(t, reply, "PONG")

	send(conn, "USE_GROUP group2")
	reply = readReply(buf)
	assert.Equal(t, reply[0:3], "Err")

	send(conn, "USE_GROUP group1")
	reply = readReply(buf)
	assert.Equal(t, reply, "OK")

	for !secondClientReady {
	}
	secondClientReady = false

	fakeJobArgs = getFakeJobArgs()
	jsonRaw, _ := json.Marshal(fakeJobArgs)
	jsonStr := string(jsonRaw)
	send(conn, "PUBLISH default job.type '"+jsonStr+"'")
	reply = readReply(buf)
	assert.Equal(t, reply, "OK")
	debug := getDebug(conn, buf)
	assert.Equal(t, debug.Queues[0].Slug, "default")
	assert.Equal(t, debug.Queues[0].Jobs[0].Type, "job.type")
	assert.Equal(t, debug.Queues[0].Jobs[0].Payload, jsonStr)

	for !secondClientReady {
	}

	send(conn, "PING")
	reply = readReply(buf)
	assert.Equal(t, reply, "PONG")
}
