package kv

import (
	"testing"
)

func TestNewKVStore(t *testing.T) {
	kvStore := NewKVStore()

	if kvStore == nil {
		t.Errorf("KVStore is nil")
	}
}

func TestSetAndGet(t *testing.T) {
	kvStore := NewKVStore()

	kvStore.Set("key", "value")

	value, err := kvStore.Get("key")

	if err != nil {
		t.Errorf("Error getting value")
	}

	if value != "value" {
		t.Errorf("Invalid value")
	}
}

func TestSetAndGetInexistent(t *testing.T) {
	kvStore := NewKVStore()

	value, _ := kvStore.Get("key")

	if value != "(nil)" {
		t.Errorf("Invalid value")
	}
}
