package store

import "time"

// StartTimestamp - The unix timestamp when the process started
var StartTimestamp int64

// StartUptimeTimer -
func StartUptimeTimer() {
	StartTimestamp = time.Now().Unix()
}

// GetUptime -
func GetUptime() int64 {
	return (time.Now().Unix() - StartTimestamp)
}
