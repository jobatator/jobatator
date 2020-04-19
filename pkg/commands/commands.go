package commands

import (
	"net"
)

// CmdInterface -
type CmdInterface struct {
	Parts map[int]string
	Conn  net.Conn
}
