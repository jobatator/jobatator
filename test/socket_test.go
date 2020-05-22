package test

import (
	"bufio"
	"encoding/json"
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestConnexion(t *testing.T) {
	startServer()

	conn := getConn()
	buf := bufio.NewReader(conn)

	// check if the connexion is alive
	send(conn, "PING")
	reply := readReply(buf)
	assert.Equal(t, reply, "PONG")

	longStr := ""
	for i := 0; i < 2000; i++ {
		longStr += "a"
	}

	// check if this can received and send long strings
	send(conn, "DEBUG_PARTS firstPart secondPart "+longStr)
	reply = readReply(buf)
	var parts []string
	err := json.Unmarshal([]byte(reply), &parts)
	assert.Equal(t, nil, err)
	assert.Equal(t, len(parts), 4)
	assert.Equal(t, len(parts[3]), len(longStr))
}
