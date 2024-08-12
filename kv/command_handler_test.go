package kv

import (
	"testing"
)

func TestHandleCommandGivenQuit(t *testing.T) {
	command, _ := HandleCommand("QUIT")

	expected := KVOperation{
		Action: COMMAND_TERMINATE_CONN,
		Mutate: false,
	}

	if !command.Equals(expected) {
		t.Errorf("Invalid command %s", (&command).String())
	}
}

func TestHandleCommandGivenSetHappyPath(t *testing.T) {
	command, _ := HandleCommand("SET key value")

	expected := KVOperation{
		Action:  COMMAND_SET_KEY,
		KeyName: "key",
		Value:   "value",
		Mutate:  true,
	}

	if !command.Equals(expected) {
		t.Errorf("Invalid command %s", (&command).String())
	}
}
