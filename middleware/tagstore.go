package middleware

import (
	"github.com/docker/distribution"
	"github.com/docker/distribution/context"
)

type tagStore struct {
	ctx   context.Context
	repo  distribution.Repository
	store Store
}

func (t *tagStore) Get(ctx context.Context, tag string) (distribution.Descriptor, error) {
	val, err := t.store.GetTag(ctx, tag)
	if err == ErrNotFound {
		return distribution.Descriptor{}, distribution.ErrTagUnknown{tag}
	}
	if err != nil {
		return distribution.Descriptor{}, err
	}
	return val, nil
}

// Tag associates the tag with the provided descriptor, updating the
// current association, if needed.
func (t *tagStore) Tag(ctx context.Context, tag string, desc distribution.Descriptor) error {
	return t.store.SetTag(ctx, tag, desc)
}

// Untag removes the given tag association
func (t *tagStore) Untag(ctx context.Context, tag string) error {
	return t.store.DeleteTag(ctx, tag)
}

// All returns the set of tags for the parent repository, as
// defined in tagStore.repo
func (t *tagStore) All(ctx context.Context) ([]string, error) {
	// TODO: Make KV a tree that allows us to store mucho data
	return []string{}, nil
}

// Lookup returns the set of tags referencing the given digest.
func (t *tagStore) Lookup(ctx context.Context, digest distribution.Descriptor) ([]string, error) {
	// TODO: treeeeez
	return []string{}, nil
}
