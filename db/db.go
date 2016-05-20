package db

import (
	"fmt"
)

// i cheated

var DB db

type db struct {
	Manifests KV
	Tags      KV
}

func init() {
	DB = db{
		Manifests: make(KV),
		Tags:      make(KV),
	}
}

var ErrNotFound = fmt.Errorf("item not found")

type KV map[string][]byte

func (store KV) Get(key string) ([]byte, error) {
	val, ok := store[key]
	if !ok {
		return nil, ErrNotFound
	}
	return val, nil
}

func (store KV) Put(key string, data []byte) error {
	store[key] = data
	return nil
}

func (store KV) Delete(key string) error {
	delete(store, key)
	return nil
}
