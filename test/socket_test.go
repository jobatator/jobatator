package test

import (
	"bufio"
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

	var parts []string

	// check if the command parser can handle quoted string
	send(conn, "DEBUG_PARTS nonQuoted 'that is quoted' 'also that' 'and' that")
	err := readJSONReply(buf, &parts)
	assert.Equal(t, nil, err)
	assert.Equal(t, len(parts), 6)
	assert.Equal(t, parts[1], "nonQuoted")
	assert.Equal(t, parts[3], "also that")
	assert.Equal(t, parts[4], "and")
	assert.Equal(t, parts[5], "that")

	// check if this can receive and send long strings
	send(conn, "DEBUG_PARTS firstPart secondPart "+longStr)
	err = readJSONReply(buf, &parts)
	assert.Equal(t, nil, err)
	assert.Equal(t, len(parts), 4)
	assert.Equal(t, len(parts[3]), len(longStr))
}
