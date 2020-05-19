package commands

import (
	"net"

	"github.com/lefuturiste/jobatator/pkg/store"
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
	cmd.Conn.Write([]byte(data))
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
