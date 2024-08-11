package main

import (
	"errors"
	"log"
	"regexp"
	"strings"
)

type CommandType string

// Define the possible values for ActionType
const (
	COMMAND_TERMINATE_CONN CommandType = "terminate-conn"
	COMMAND_RETURN_KEY     CommandType = "return-key"
)

type KVOperation struct {
	Action  CommandType
	KeyName string
	replyCh chan string
}

func handleQuit(_arg1 string, _arg2 string) KVOperation {
	return KVOperation{
		Action: COMMAND_TERMINATE_CONN,
	}
}

func handleGet(keyName string, _arg2 string) KVOperation {
	return KVOperation{
		Action:  COMMAND_RETURN_KEY,
		KeyName: keyName,
	}
}

func HandleCommand(rawCommand string) (KVOperation, error) {
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
		arg1 := ""
		arg2 := ""

		if len(matches) > 3 {
			arg1 = matches[3]
		}

		if len(matches) > 5 {
			arg2 = matches[5]
		}

		commands := map[string]func(string, string) KVOperation{
			"QUIT": handleQuit,
			"GET":  handleGet,
		}

		if commandFunc, found := commands[commandType]; found {
			return commandFunc(arg1, arg2), nil
		} else {
			log.Printf("Unknown command type: %s", commandType)
		}

	}

	return KVOperation{}, errors.New("No valid command found")
}
