package middleware

import (
	"encoding/json"
	"fmt"

	"tagstore/db"

	"github.com/docker/distribution"
	"github.com/docker/distribution/context"
	"github.com/docker/distribution/digest"
	"github.com/docker/distribution/manifest"
	"github.com/docker/distribution/manifest/manifestlist"
	"github.com/docker/distribution/manifest/schema1"
	"github.com/docker/distribution/manifest/schema2"
)

// manifestStore provides an alternative backing mechanism for manifests.
// It must implement the ManifestService to store manifests and
// ManifestEnumerator for garbage collection and listing
type manifestStore struct {
	ctx        context.Context
	repository distribution.Repository
}

func (m *manifestStore) Exists(ctx context.Context, dgst digest.Digest) (bool, error) {
	_, err := db.DB.Manifests.Get(string(dgst))
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
	content, err := db.DB.Manifests.Get(string(dgst))
	if err != nil {
		return nil, err
	}

	var versioned manifest.Versioned
	if err = json.Unmarshal(content, &versioned); err != nil {
		return nil, err
	}

	switch versioned.SchemaVersion {
	case 1:
		return m.UnmarshalV1(ctx, dgst, content)
	case 2:
		// This can be an image manifest or a manifest list
		switch versioned.MediaType {
		case schema2.MediaTypeManifest:
			return m.UnmarshalV2(ctx, dgst, content)
		case manifestlist.MediaTypeManifestList:
			return m.UnmarshalList(ctx, dgst, content)
		default:
			return nil, distribution.ErrManifestVerification{fmt.Errorf("unrecognized manifest content type %s", versioned.MediaType)}
		}
	}

	return nil, distribution.ErrManifestBlobUnknown{dgst}
}

// Put creates or updates the given manifest returning the manifest digest
func (m *manifestStore) Put(ctx context.Context, manifest distribution.Manifest, options ...distribution.ManifestServiceOption) (d digest.Digest, err error) {

	// NOTE: we're not allowing skipDependencyVerification here.
	//
	// skipDependencyVerification is ONLY used when registry is set up as a
	// pull-through cache (proxy). In these circumstances this middleware
	// should not be used, therefore this verification implementation always
	// verifies blobs.
	//
	// This is the only difference in implementation with storage's
	// manifestStore{}
	switch manifest.(type) {
	case *schema1.SignedManifest:
		err = m.VerifyV1(ctx, manifest.(*schema1.SignedManifest))
	case *schema2.DeserializedManifest:
		err = m.VerifyV2(ctx, manifest.(*schema2.DeserializedManifest))
	case *manifestlist.DeserializedManifestList:
		err = m.VerifyList(ctx, manifest.(*manifestlist.DeserializedManifestList))
	default:
		err = fmt.Errorf("Unknown manifest type: %T", manifest)
	}

	if err != nil {
		return
	}

	// The database needs only store the bytes here; we'll decode into
	// manifest.Versioned in order to detect the version.
	_, data, err := manifest.Payload()
	if err != nil {
		return
	}

	dgst := digest.FromBytes(data)
	err = db.DB.Manifests.Put(string(dgst), data)

	return dgst, err
}

// Delete removes the manifest specified by the given digest
func (m *manifestStore) Delete(ctx context.Context, dgst digest.Digest) error {
	if _, err := db.DB.Manifests.Get(string(dgst)); err != nil {
		return err
	}
	return db.DB.Manifests.Delete(string(dgst))
}
