package commands

import (
	"encoding/json"
	"runtime"

	"github.com/lefuturiste/jobatator/pkg/store"
	"github.com/lefuturiste/jobatator/pkg/utils"
)

// HealthOutput -
type HealthOutput struct {
	Uptime                   int64
	Os                       string
	Version                  string
	NumCPU                   int
	NumGoroutine             int
	NumSessions              int
	NumAuthenticatedSessions int
}

// Health -
func Health(cmd CmdInterface) {
	var output HealthOutput
	output.Uptime = utils.GetUptime()
	output.Os = runtime.GOOS
	output.Version = runtime.Version()
	output.NumCPU = runtime.NumCPU()
	output.NumGoroutine = runtime.NumGoroutine()
	output.NumSessions = len(store.Sessions)
	for _, session := range store.Sessions {
		if session.Username != "" {
			output.NumAuthenticatedSessions++
		}
	}
	rawJSON, _ := json.Marshal(output)
	ReturnString(cmd, string(rawJSON))
}
