package middleware

import (
	"encoding/json"

	"github.com/docker/distribution"
	"github.com/docker/distribution/context"
	"github.com/docker/distribution/digest"
	"github.com/docker/distribution/manifest/schema2"
)

func (m *manifestStore) VerifyV2(ctx context.Context, mnfst *schema2.DeserializedManifest) error {
	var errs distribution.ErrManifestVerification

	target := mnfst.Target()
	_, err := m.repository.Blobs(ctx).Stat(ctx, target.Digest)
	if err != nil {
		if err != distribution.ErrBlobUnknown {
			errs = append(errs, err)
		}

		// On error here, we always append unknown blob errors.
		errs = append(errs, distribution.ErrManifestBlobUnknown{Digest: target.Digest})
	}

	for _, fsLayer := range mnfst.References() {
		_, err := m.repository.Blobs(ctx).Stat(ctx, fsLayer.Digest)
		if err != nil {
			if err != distribution.ErrBlobUnknown {
				errs = append(errs, err)
			}

			// On error here, we always append unknown blob errors.
			errs = append(errs, distribution.ErrManifestBlobUnknown{Digest: fsLayer.Digest})
		}
	}

	if len(errs) != 0 {
		return errs
	}

	return nil
}

func (m *manifestStore) UnmarshalV2(ctx context.Context, dgst digest.Digest, content []byte) (distribution.Manifest, error) {
	context.GetLogger(m.ctx).Debug("(*schema2ManifestHandler).Unmarshal")

	var man schema2.DeserializedManifest
	if err := json.Unmarshal(content, &man); err != nil {
		return nil, err
	}

	return &man, nil
}
