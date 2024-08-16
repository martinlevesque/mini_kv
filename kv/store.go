package kv

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"
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

			log.Printf("Mutating loop - %s", currentKvOperation.Action)

			if currentKvOperation.Action == COMMAND_SET_KEY {
				kvStore.Set(currentKvOperation.KeyName, currentKvOperation.Value)
			} else if currentKvOperation.Action == COMMAND_DEL_KEY {
				kvStore.Del(currentKvOperation.KeyName)
			}

			currentKvOperation.ReplyCh <- "(nil)"
		}
	}(mutableCommandsChannel)

	return &kvStore
}

func (kvStore *KVStore) ImmutableOperation(op *KVOperation) (string, error) {
	log.Printf("Immutable operation: %s", op.Action)
	if op.Action == COMMAND_RETURN_KEY {
		value, err := kvStore.Get(op.KeyName)

		if err != nil {
			return "", err
		}

		return value, nil
	} else if op.Action == COMMAND_EXPIRE {
		go kvStore.Expire(op)
		log.Printf("after expire: %s", op.KeyName)
		return "", nil
	} else if op.Action == COMMAND_KEYS {
		return kvStore.Keys(op.KeyName)
	}

	return "", errors.New("Invalid operation")
}

func (kvStore *KVStore) Keys(regexPattern string) (string, error) {
	pattern, err := regexp.Compile(regexPattern)

	if err != nil {
		return "", err // Return the error if the pattern is invalid
	}

	keys := ""
	i := 0

	for key := range kvStore.store {
		i++

		if pattern.MatchString(key) {
			keys += fmt.Sprintf("%d) %s\n", i, key)
		}
	}

	return keys, nil
}

func (kvStore *KVStore) Expire(op *KVOperation) {
	log.Println("IN Expire not implemented yet")

	// Value to int
	expire_in, err := strconv.Atoi(op.Value)

	if err != nil {
		log.Printf("Failed to convert expire value to int: %s", err)
		return
	}

	time.Sleep(time.Duration(expire_in) * time.Second)
	log.Printf("Expired key: %s", op.KeyName)

	opDelete := KVOperation{
		Action:  COMMAND_DEL_KEY,
		KeyName: op.KeyName,
		Value:   "",
		Mutate:  true,
		ReplyCh: make(chan string, 1),
	}
	kvStore.MutableCommandsChannel <- opDelete

	// ignore the ReplyCh
	<-opDelete.ReplyCh
}

func (kvStore *KVStore) Set(key string, value string) {
	kvStore.store[key] = value
}

func (kvStore *KVStore) Del(key string) {
	delete(kvStore.store, key)
}

// Get gets a key from the store
func (kvStore *KVStore) Get(key string) (string, error) {
	value, ok := kvStore.store[key]

	if !ok {
		return "(nil)", nil
	}

	return string([]byte(value)), nil
}
