package middleware

import (
	"tagstore/db"

	"github.com/docker/distribution"
	"github.com/docker/distribution/context"
	"github.com/docker/distribution/digest"
)

type tagStore struct {
	ctx  context.Context
	repo distribution.Repository
}

func (t *tagStore) Get(ctx context.Context, tag string) (distribution.Descriptor, error) {
	val, err := db.DB.Tags.Get(tag)
	if err == db.ErrNotFound {
		return distribution.Descriptor{}, distribution.ErrTagUnknown{tag}
	}
	if err != nil {
		return distribution.Descriptor{}, err
	}

	// TODO: should tags store media type for manifests?

	return distribution.Descriptor{
		Digest: digest.Digest(val),
	}, nil
}

// Tag associates the tag with the provided descriptor, updating the
// current association, if needed.
func (t *tagStore) Tag(ctx context.Context, tag string, desc distribution.Descriptor) error {
	return db.DB.Tags.Put(tag, []byte(desc.Digest))
}

// Untag removes the given tag association
func (t *tagStore) Untag(ctx context.Context, tag string) error {
	return db.DB.Tags.Delete(tag)
}

// All returns the set of tags managed by this tag service
func (t *tagStore) All(ctx context.Context) ([]string, error) {
	// TODO: Make KV a tree that allows us to store mucho data
	return []string{}, nil
}

// Lookup returns the set of tags referencing the given digest.
func (t *tagStore) Lookup(ctx context.Context, digest distribution.Descriptor) ([]string, error) {
	// TODO: treeeeez
	return []string{}, nil
}
