package main

import (
	"errors"
	"testing"
)

func TestPut(t *testing.T) {
	key := "create-key"
	value := "create-value"

	var val interface{}
	var contains bool

	defer delete(store, key)

	_, contains = store[key]
	if contains {
		t.Error("k/v already exists")
	}

	err := Put(key, value)
	if err != nil {
		t.Error(err)
	}

	val, contains = store[key]
	if !contains {
		t.Error("create failed")
	}

	if val != value {
		t.Error("val/value mismatch")
	}
}

func TestGet(t *testing.T) {
	key := "read-key"
	value := "read-value"

	var val interface{}
	var err error

	defer delete(store, key)

	val, err = Get(key)
	if err == nil {
		t.Error("expected error:", err)
	}

	if !errors.Is(err, ErrorNoSuchKey) {
		t.Error("unexpected error:", err)
	}

	store[key] = value

	val, err = Get(key)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	if val != value {
		t.Error("val/value mismatch")
	}
}

func TestDelete(t *testing.T) {
	key := "delete-key"
	value := "delete-value"

	var contains bool

	defer delete(store, key)

	store[key] = value

	_, contains = store[key]
	if !contains {
		t.Error("k/v do not exist")
	}

	Delete(key)

	_, contains = store[key]
	if contains {
		t.Error("Delete failed")
	}
}
