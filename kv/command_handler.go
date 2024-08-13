package kv

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
)

type CommandType string

// Define the possible values for ActionType
const (
	COMMAND_TERMINATE_CONN CommandType = "terminate-conn"
	COMMAND_RETURN_KEY     CommandType = "return-key"
	COMMAND_SET_KEY        CommandType = "set-key"
)

type KVOperation struct {
	Action  CommandType
	KeyName string
	Value   string
	ReplyCh chan string
	Mutate  bool
}

func (kv *KVOperation) String() string {
	return fmt.Sprintf("KVOperation{Action: %s, KeyName: %s, Value: %s, Mutate: %t}",
		kv.Action, kv.KeyName, kv.Value, kv.Mutate)
}

func (kv KVOperation) Equals(other KVOperation) bool {
	return kv.Action == other.Action &&
		kv.KeyName == other.KeyName &&
		kv.Value == other.Value &&
		kv.Mutate == other.Mutate
}

func handleQuit(_arg1 string, _arg2 string) KVOperation {
	return KVOperation{
		Action: COMMAND_TERMINATE_CONN,
		Mutate: false,
	}
}

func handleSet(keyName string, value string) KVOperation {
	return KVOperation{
		Action:  COMMAND_SET_KEY,
		KeyName: keyName,
		Value:   value,
		Mutate:  true,
	}
}

func handleGet(keyName string, _arg2 string) KVOperation {
	return KVOperation{
		Action:  COMMAND_RETURN_KEY,
		KeyName: keyName,
		Mutate:  false,
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
			"SET":  handleSet,
		}

		if commandFunc, found := commands[commandType]; found {
			result := commandFunc(arg1, arg2)

			result.ReplyCh = make(chan string)

			return result, nil
		} else {
			log.Printf("Unknown command type: %s", commandType)
		}

	}

	return KVOperation{}, errors.New("No valid command found")
}
