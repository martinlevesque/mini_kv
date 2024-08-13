package kv

import (
	"errors"
	"log"
)

// KVStore is a simple key-value store in memory

type KVStore struct {
	store                  map[string]string
	MutableCommandsChannel chan KVOperation
}

// NewKVStore creates a new KVStore
func NewKVStore() *KVStore {
	log.Println("Creating new KVStore")

	mutableCommandsChannel := make(chan KVOperation)

	kvStore := KVStore{
		store:                  make(map[string]string),
		MutableCommandsChannel: mutableCommandsChannel,
	}

	go func(kvOperation <-chan KVOperation) {
		for {
			// Wait for a command from the channel
			currentKvOperation := <-kvOperation

			if currentKvOperation.Action == COMMAND_SET_KEY {
				kvStore.Set(currentKvOperation.KeyName, currentKvOperation.Value)
			}

			currentKvOperation.ReplyCh <- "(nil)"
		}
	}(mutableCommandsChannel)

	return &kvStore
}

func (kvStore *KVStore) ImmutableOperation(op *KVOperation) (string, error) {
	if op.Action == COMMAND_RETURN_KEY {
		value, err := kvStore.Get(op.KeyName)

		if err != nil {
			return "", err
		}

		return value, nil
	}

	return "", errors.New("Invalid operation")
}

// Set sets a key in the store
func (kvStore *KVStore) Set(key string, value string) {
	kvStore.store[key] = value
}

// Get gets a key from the store
func (kvStore *KVStore) Get(key string) (string, error) {
	value, ok := kvStore.store[key]

	if !ok {
		return "(nil)", nil
	}

	return string([]byte(value)), nil
}
