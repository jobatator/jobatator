package test

import (
	"bufio"
	"encoding/json"
	"testing"

	"github.com/lefuturiste/jobatator/pkg/server"
	"github.com/lefuturiste/jobatator/pkg/store"
	"github.com/magiconair/properties/assert"
)

const simpleConfig = `
port: 8963
host: "0.0.0.0"
test_mode: true
log_level: "DEBUG"
`

func TestConnexion(t *testing.T) {
	store.LoadConfigFromString(simpleConfig)
	server.StartAsync()

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
