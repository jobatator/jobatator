package server

import "github.com/lefuturiste/jobatator/pkg/commands"

// CmdMap -
var CmdMap = map[string]interface{}{
	// connexion
	"PING":        commands.Ping,
	"QUIT":        commands.Quit,
	"AUTH":        commands.Auth,
	"USE_GROUP":   commands.UseGroup,
	"PUBLISH":     commands.Publish,
	"SUBSCRIBE":   commands.Subscribe,
	"UPDATE_JOB":  commands.UpdateJob,
	"DISPATCH":    commands.Dispatch,
	"DEBUG":       commands.Debug,
	"DEBUG_JSON":  commands.DebugJSON,
	"STOP_SERVER": commands.StopServer,
}
