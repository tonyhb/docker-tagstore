package db

import (
	"tagstore/middleware"

	"github.com/docker/distribution"
	"github.com/docker/distribution/context"
)

type ManifestKV map[string][]byte

func (store ManifestKV) SetManifest(ctx context.Context, key string, val distribution.Manifest) error {
	_, byt, err := val.Payload()
	if err != nil {
		return err
	}
	store[key] = byt
	return nil
}

func (store ManifestKV) GetManifest(ctx context.Context, key string) ([]byte, error) {
	val, ok := store[key]
	if !ok {
		return nil, middleware.ErrNotFound
	}
	return val, nil
}

func (store ManifestKV) DeleteManifest(ctx context.Context, key string) error {
	if _, ok := store[key]; !ok {
		return middleware.ErrNotFound
	}
	delete(store, key)
	return nil
}
