package commands

import (
	"net"

	"github.com/jobatator/jobatator/pkg/store"
)

// CmdDefinition - Use to define the command
type CmdDefinition struct {
	Name        string
	RequireAuth bool
	UseGroup    bool
	Args        int
	Description string
	Usage       string
	Handler     func(CmdInterface)
}

// CmdInterface - Passed to command handler
type CmdInterface struct {
	Parts map[int]string
	Conn  net.Conn
	User  store.User
}

// ReturnString -
func ReturnString(cmd CmdInterface, data string) {
	var buf []byte
	for _, char := range []byte(data) {
		if len(buf) == 1024 {
			cmd.Conn.Write(buf)
			buf = []byte{}
		}
		buf = append(buf, char)
	}
	if len(buf) > 0 {
		cmd.Conn.Write(buf)
	}
	NewLine(cmd)
}

// ReturnError -
func ReturnError(cmd CmdInterface, err string) {
	ReturnString(cmd, "Err: "+err)
}

// NewLine -
func NewLine(cmd CmdInterface) {
	cmd.Conn.Write([]byte("\n"))
}
