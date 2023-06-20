package main

import "errors"

var store = make(map[string]string)
var ErrorNoSuchKey = errors.New("no such key")

func Put(key, value string) error {
	store[key] = value

	return nil
}

func Get(key string) (string, error) {
	if _, exists := store[key]; exists {
		return store[key], nil
	}
	return "", ErrorNoSuchKey
}

func Delete(key string) error {
	delete(store, key)
	return nil
}
