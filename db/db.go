package db

import (
	"github.com/docker/distribution"
	"github.com/docker/distribution/digest"
)

// i cheated

var manifests = map[digest.Digest][]byte{}

func init() {
	manifests = make(map[digest.Digest][]byte)
}

func Get(dgst digest.Digest) ([]byte, error) {
	val, ok := manifests[dgst]
	if !ok {
		return nil, distribution.ErrManifestBlobUnknown{dgst}
	}
	return val, nil
}

func Put(data []byte) (digest.Digest, error) {
	dgst := digest.FromBytes(data)
	manifests[dgst] = data
	return dgst, nil
}

func Delete(dgst digest.Digest) error {
	delete(manifests, dgst)
	return nil
}
