package db

import (
	"fmt"

	"github.com/docker/distribution"
	"github.com/docker/distribution/digest"
)

// i cheated

var manifests = map[digest.Digest]interface{}{}

func init() {
	manifests = make(map[digest.Digest]interface{})
}

func Get(dgst digest.Digest) (distribution.Manifest, error) {
	val, ok := manifests[dgst]
	if !ok {
		return nil, distribution.ErrManifestBlobUnknown{dgst}
	}
	man, ok := val.(distribution.Manifest)
	if !ok {
		return nil, fmt.Errorf("manifest not found: %s", dgst)
	}
	return man, nil
}

func Put(manifest distribution.Manifest) (digest.Digest, error) {
	_, hash, err := manifest.Payload()
	if err != nil {
		return "", err
	}

	dgst := digest.FromBytes(hash)
	manifests[dgst] = manifest

	return dgst, nil
}

func Delete(dgst digest.Digest) error {
	delete(manifests, dgst)
	return nil
}
