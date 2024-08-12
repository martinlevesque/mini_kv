package kv

import (
	"errors"
	"log"
)

// KVStore is a simple key-value store in memory

type KVStore struct {
	store map[string]string
}

// NewKVStore creates a new KVStore
func NewKVStore() *KVStore {
	log.Println("Creating new KVStore")

	return &KVStore{
		store: make(map[string]string),
	}
}

// Set sets a key in the store
func (kv *KVStore) Set(key string, value string) {
	kv.store[key] = value
}

// Get gets a key from the store
func (kv *KVStore) Get(key string) (string, error) {
	//return string([]byte(kv.store[key]))
	value, ok := kv.store[key]

	if !ok {
		return "", errors.New("Key not found")
	}

	return string([]byte(value)), nil
}
