package test

import (
	"bufio"
	"encoding/json"
	"net"
	"strings"

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
