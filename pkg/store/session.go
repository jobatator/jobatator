package store

import "net"

// Sessions - all the sessions
var Sessions []User

// WorkerBusy -
const WorkerBusy = "busy"

// WorkerAvailable -
const WorkerAvailable = "available"

// FindSession - will return a user object for this session
func FindSession(conn net.Conn) User {
	var user User
	for _, value := range Sessions {
		if value.Addr == conn.RemoteAddr().String() {
			user = value
		}
	}
	return user
}

// UpdateSession -
func UpdateSession(user User) {
	for key, value := range Sessions {
		if value.Addr == user.Addr {
			Sessions[key] = user
		}
	}
}
