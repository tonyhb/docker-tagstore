package middleware

import (
	"tagstore/db"

	"github.com/docker/distribution"
	"github.com/docker/distribution/context"
	"github.com/docker/distribution/digest"
)

// manifestStore provides an alternative backing mechanism for manifests.
// It must implement the ManifestService to store manifests and
// ManifestEnumerator for garbage collection and listing
type manifestStore struct {
	ctx        context.Context
	repository distribution.Repository

	deleteEnabled bool
}

func (m *manifestStore) Exists(ctx context.Context, dgst digest.Digest) (bool, error) {
	_, err := db.Get(dgst)
	if _, ok := err.(distribution.ErrManifestBlobUnknown); ok {
		// TODO: return an ErrManifestUnknownRevision
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// Get retrieves the manifest specified by the given digest.
// Note that the middleware itself verifies that the manifest is valid;
// the storage backend should only marshal and unmarshal into the correct type.
func (m *manifestStore) Get(ctx context.Context, dgst digest.Digest, options ...distribution.ManifestServiceOption) (distribution.Manifest, error) {
	return db.Get(dgst)
}

// Put creates or updates the given manifest returning the manifest digest
func (m *manifestStore) Put(ctx context.Context, manifest distribution.Manifest, options ...distribution.ManifestServiceOption) (digest.Digest, error) {
	return db.Put(manifest)
}

// Delete removes the manifest specified by the given digest
func (m *manifestStore) Delete(ctx context.Context, dgst digest.Digest) error {
	if _, err := db.Get(dgst); err != nil {
		return err
	}
	return db.Delete(dgst)
}
