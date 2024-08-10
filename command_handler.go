package main

import (
	"log"
	"regexp"
	"strings"
)

func HandleCommand(rawCommand string) {
	command := strings.TrimSpace(rawCommand)
	log.Printf("Received command: %s", command)

	// QUIT
	// GET <key>
	// SET <key> <value>
	// DEL <key>
	// KEYS
	// EXPIRE <key> <seconds>

	re := regexp.MustCompile(`^(QUIT|GET|SET|DEL|KEYS|EXPIRE)(\s+([^\s]+))?(\s+([^\s]+))?$`)
	matches := re.FindStringSubmatch(command)

	if len(matches) > 1 {
		// First group will be the command
		commandType := matches[1]
		log.Printf("Command type: %s", commandType)

		if len(matches) > 2 {
			log.Printf("Command argument group 2: %s", matches[2])
		}
	} else {
		log.Printf("No valid command found")
	}

}
