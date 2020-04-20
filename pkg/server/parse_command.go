package server

import "strings"

func parseCommand(input string) map[int]string {
	//var input = `SET hello-world '{"random": true}'`
	var cmdComponents = strings.Split(input, " ")
	var quotedComponent = ""
	var inQuoted bool = false
	var components map[int]string = make(map[int]string)
	var componentIndex int = 0
	for _, value := range cmdComponents {
		if strings.Contains(value, "'") && value[0:1] == "'" {
			// we encountered a new quoted component
			// fmt.Println("Start of quoted ", value)
			if value[len(value)-1:] == "'" {
				components[componentIndex] = value[1 : len(value)-1]
				componentIndex++
			} else {
				inQuoted = true
				quotedComponent = value
			}
		} else if strings.Contains(value, "'") && value[len(value)-1:] == "'" {
			// fmt.Println("End of quoted ", value)
			// we reached the end of the quoted component
			quotedComponent += " " + value
			quotedComponent = quotedComponent[1 : len(quotedComponent)-1]
			inQuoted = false
			components[componentIndex] = quotedComponent
			componentIndex++
		} else if inQuoted {
			quotedComponent += " " + value
		} else {
			components[componentIndex] = value
			componentIndex++
		}
	}

	return components
}
