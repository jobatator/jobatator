package commands

import (
	"encoding/json"
	"runtime"

	"github.com/lefuturiste/jobatator/pkg/store"
	"github.com/lefuturiste/jobatator/pkg/utils"
)

// HealthOutput -
type HealthOutput struct {
	Uptime        int64
	Os            string
	Version       string
	NumCPU        int
	NumGoroutine  int
	NumSessions   int
	NumAuthed     int
	NumGroups     int
	NumUsers      int
	NumQueues     int
	NumJobs       int
	NumWorkers    int
	UnixTimestamp int64
}

// Health - Get all kinds of data about the jobatator instance
func Health(cmd CmdInterface) {
	var output HealthOutput
	output.Uptime = utils.GetUptime()
	output.Os = runtime.GOOS
	output.Version = runtime.Version()
	output.NumCPU = runtime.NumCPU()
	output.NumGoroutine = runtime.NumGoroutine()
	output.UnixTimestamp = utils.GetUnixTimestamp()

	// list sessions
	output.NumSessions = len(store.Sessions)
	for _, session := range store.Sessions {
		if session.Username != "" {
			output.NumAuthed++
		}
	}

	// list queues, jobs and workers
	output.NumQueues = len(store.Queues)
	output.NumJobs = 0
	for _, queue := range store.Queues {
		output.NumJobs += len(queue.Jobs)
		output.NumWorkers += len(queue.Workers)
	}

	rawJSON, _ := json.Marshal(output)
	ReturnString(cmd, string(rawJSON))
}
