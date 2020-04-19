package server

import "github.com/lefuturiste/jobatator/pkg/commands"

// CmdMap -
var CmdMap = map[string]interface{}{
	// connexion
	"PING":      commands.Ping,
	"QUIT":      commands.Quit,
	"AUTH":      commands.Auth,
	"USE_GROUP": commands.UseGroup,
	"PUBLISH":   commands.Publish,
	"SUBSCRIBE": commands.Subscribe,
	"DEBUG":     commands.Debug,
}
