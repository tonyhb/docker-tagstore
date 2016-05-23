package db

import (
	"tagstore/middleware"

	"github.com/docker/distribution"
	"github.com/docker/distribution/context"
	"github.com/docker/distribution/digest"
)

type TagKV map[string][]byte

func (store TagKV) GetTag(ctx context.Context, key string) (distribution.Descriptor, error) {
	val, ok := store[key]
	if !ok {
		return distribution.Descriptor{}, middleware.ErrNotFound
	}
	dgst, err := digest.ParseDigest(string(val))
	if err != nil {
		return distribution.Descriptor{}, err
	}
	return distribution.Descriptor{Digest: dgst}, nil
}

func (store TagKV) SetTag(ctx context.Context, key string, val distribution.Descriptor) error {
	store[key] = []byte(val.Digest)
	return nil
}

func (store TagKV) DeleteTag(ctx context.Context, key string) error {
	if _, ok := store[key]; !ok {
		return middleware.ErrNotFound
	}
	delete(store, key)
	return nil
}

func (store TagKV) AllTags(ctx context.Context, repo string) ([]string, error) {
	// TODO
	return []string{}, nil
}

func (store TagKV) Lookup(ctx context.Context, digest distribution.Descriptor) ([]string, error) {
	// TODO
	return []string{}, nil
}
