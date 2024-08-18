package kv

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"sync"
	"time"
)

// KVStore is a simple key-value store in memory

type KVStore struct {
	store                  map[string]string
	Mutex                  sync.RWMutex
	MutableCommandsChannel chan KVOperation
}

// NewKVStore creates a new KVStore
func NewKVStore() *KVStore {
	log.Println("Creating new KVStore")

	mutableCommandsChannel := make(chan KVOperation)

	kvStore := KVStore{
		store:                  make(map[string]string),
		MutableCommandsChannel: mutableCommandsChannel,
		Mutex:                  sync.RWMutex{},
	}

	go func(kvOperation <-chan KVOperation) {
		for {
			// Wait for a command from the channel
			currentKvOperation := <-kvOperation

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
	if op.Action == COMMAND_RETURN_KEY {
		value, err := kvStore.Get(op.KeyName)

		if err != nil {
			return "", err
		}

		return value, nil
	} else if op.Action == COMMAND_EXPIRE {
		go kvStore.Expire(op)
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

	kvStore.Mutex.RLock()
	for key := range kvStore.store {
		i++

		if pattern.MatchString(key) {
			keys += fmt.Sprintf("%d) %s\n", i, key)
		}
	}
	kvStore.Mutex.RUnlock()

	keys += fmt.Sprintf("Total keys: %d\n", i)

	return keys, nil
}

func (kvStore *KVStore) Expire(op *KVOperation) {
	// Value to int
	expire_in, err := strconv.Atoi(op.Value)

	if err != nil {
		log.Printf("Failed to convert expire value to int: %s", err)
		return
	}

	time.Sleep(time.Duration(expire_in) * time.Second)

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
	kvStore.Mutex.Lock()
	kvStore.store[key] = value
	kvStore.Mutex.Unlock()
}

func (kvStore *KVStore) Del(key string) {
	kvStore.Mutex.Lock()
	delete(kvStore.store, key)
	kvStore.Mutex.Unlock()
}

func (kvStore *KVStore) Get(key string) (string, error) {
	kvStore.Mutex.RLock()
	value, ok := kvStore.store[key]
	kvStore.Mutex.RUnlock()

	if !ok {
		return "(nil)", nil
	}

	return string([]byte(value)), nil
}
