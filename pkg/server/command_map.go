package server

import (
	cmds "github.com/lefuturiste/jobatator/pkg/commands"
)

// CmdMap -
var CmdMap = []cmds.CmdDefinition{
	{
		Name:        "PING",
		Handler:     cmds.Ping,
		RequireAuth: false,
		UseGroup:    false,
	}, {
		Name:        "QUIT",
		Handler:     cmds.Quit,
		RequireAuth: false,
		UseGroup:    false,
	}, {
		Name:        "AUTH",
		Handler:     cmds.Auth,
		RequireAuth: false,
		UseGroup:    false,
		Args:        2,
	}, {
		Name:        "USE_GROUP",
		Handler:     cmds.UseGroup,
		RequireAuth: true,
		UseGroup:    false,
		Args:        1,
	}, {
		Name:        "DEBUG",
		Handler:     cmds.Debug,
		RequireAuth: false,
		UseGroup:    false,
	}, {
		Name:        "DEBUG_JSON",
		Handler:     cmds.DebugJSON,
		RequireAuth: false,
		UseGroup:    false,
	}, {
		Name:        "STOP_SERVER",
		Handler:     cmds.StopServer,
		RequireAuth: false,
		UseGroup:    false,
	}, {
		Name:        "PUBLISH",
		Handler:     cmds.Publish,
		RequireAuth: true,
		UseGroup:    true,
		Args:        3,
	}, {
		Name:        "SUBSCRIBE",
		Handler:     cmds.Subscribe,
		RequireAuth: true,
		UseGroup:    true,
		Args:        1,
	}, {
		Name:        "UPDATE_JOB",
		Handler:     cmds.UpdateJob,
		RequireAuth: true,
		UseGroup:    true,
		Args:        2,
	}, {
		Name:        "DELETE_JOB",
		Handler:     cmds.DeleteJob,
		RequireAuth: true,
		UseGroup:    true,
		Args:        1,
	}, {
		Name:        "DISPATCH",
		Handler:     cmds.Dispatch,
		RequireAuth: true,
		UseGroup:    true,
	},
}
