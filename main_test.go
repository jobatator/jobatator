package main

import (
	"bufio"
	"net"
	"strings"
	"testing"

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

func TestJob(t *testing.T) {
	utils.LoadConfigFromString(testConfig)
	server.StartAsync()

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

}
