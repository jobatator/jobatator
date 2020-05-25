package test

import (
	"bufio"
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestAuth(t *testing.T) {
	startServer()
	conn := getConn()
	doAuthStuff(conn)
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

	// authenticate with super user
	send(conn, "AUTH superuser superuser")
	reply = readReply(buf)
	assert.Equal(t, reply, "Welcome!")

	// try to access a group with the superuser
	send(conn, "USE_GROUP group1")
	reply = readReply(buf)
	assert.Equal(t, reply, "OK")
}
