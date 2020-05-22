package test

import (
	"bufio"
	"net"
	"strings"
)

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
