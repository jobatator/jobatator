package commands

import (
	"encoding/json"
	"fmt"
)

// Help - List all commands or get help about a specific command
func Help(cmd CmdInterface) {
	if len(cmd.Parts) == 2 {
		// get help about a specific command
		var command CmdDefinition
		for _, value := range CmdMap {
			if value.Name == cmd.Parts[1] {
				command = value
			}
		}
		if command.Name == "" {
			ReturnError(cmd, "unknown-command")
			return
		}
		type CmdHelpOutput struct {
			Name         string
			RequireAuth  bool
			RequireGroup bool
			Args         int
			Description  string
			Usage        string
		}
		cmdOutput := CmdHelpOutput{
			Name:         command.Name,
			Description:  command.Description,
			Usage:        command.Usage,
			Args:         command.Args,
			RequireAuth:  command.RequireAuth,
			RequireGroup: command.UseGroup,
		}
		if cmdOutput.Usage == "" {
			cmdOutput.Usage = cmdOutput.Name
		}
		rawJSON, _ := json.Marshal(cmdOutput)
		fmt.Println(string(rawJSON))
		ReturnString(cmd, string(rawJSON))
	} else {
		// list all commands
		var commands []string
		for _, value := range CmdMap {
			commands = append(commands, value.Name)
		}
		rawJSON, _ := json.Marshal(commands)
		ReturnString(cmd, string(rawJSON))
	}
}
