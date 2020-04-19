package utils

import (
	"net"
)

// CmdInterface -
type CmdInterface struct {
	Parts map[int]string
	Conn  net.Conn
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
