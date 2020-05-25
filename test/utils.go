package test

import (
	"bufio"
	"encoding/json"
	"net"
	"strconv"
	"strings"

	"github.com/dchest/uniuri"
	"github.com/lefuturiste/jobatator/pkg/commands"
	"github.com/lefuturiste/jobatator/pkg/server"
	"github.com/lefuturiste/jobatator/pkg/store"
)

const testConfig = `
port: 8964
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
    - username: "superuser"
      password: "superuser"
      groups: ["*"]

delay_policy: "IGNORE"
test_mode: true
log_level: "DEBUG"
job_timeout: 2
`

var serverStarted bool

func startServer() {
	if !serverStarted {
		store.LoadConfigFromString(testConfig)
		server.StartAsync()
		serverStarted = true
	}
}

func doAuthStuff(conn net.Conn) {
	buf := bufio.NewReader(conn)
	send(conn, "AUTH user1 pass1")
	readReply(buf)
	send(conn, "USE_GROUP group1")
	readReply(buf)
}

func getConn() net.Conn {
	tcpAddr, _ := net.ResolveTCPAddr("tcp", "localhost:8964")
	conn, _ := net.DialTCP("tcp", nil, tcpAddr)
	return conn
}

func readJSONReply(buf *bufio.Reader, value interface{}) error {
	return json.Unmarshal([]byte(readReply(buf)), &value)
}

func readReply(buf *bufio.Reader) string {
	reply, _ := buf.ReadString('\n')
	return strings.Replace(reply, "\n", "", 1)
}

func send(conn net.Conn, str string) {
	conn.Write([]byte(str + "\n"))
}

func getDebug(conn net.Conn, buf *bufio.Reader) commands.DebugOutput {
	send(conn, "DEBUG_JSON")
	reply := readReply(buf)
	var debugData commands.DebugOutput
	json.Unmarshal([]byte(reply), &debugData)
	return debugData
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
