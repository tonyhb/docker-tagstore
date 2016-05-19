package middleware

import (
	"encoding/json"
	"fmt"

	"github.com/docker/distribution"
	"github.com/docker/distribution/context"
	"github.com/docker/distribution/digest"
	"github.com/docker/distribution/manifest/schema1"
	"github.com/docker/distribution/reference"
	"github.com/docker/libtrust"
)

// VerifyV1 ensures that the v1 signed manifest content is valid from the
// perspective of the registry. It ensures that the signature is valid for the
// enclosed payload. As a policy, the registry only tries to store valid
// content, leaving trust policies of that content up to consumers.
func (ms *manifestStore) VerifyV1(ctx context.Context, mnfst *schema1.SignedManifest) error {
	var errs distribution.ErrManifestVerification

	if len(mnfst.Name) > reference.NameTotalLengthMax {
		errs = append(errs,
			distribution.ErrManifestNameInvalid{
				Name:   mnfst.Name,
				Reason: fmt.Errorf("manifest name must not be more than %v characters", reference.NameTotalLengthMax),
			})
	}

	if !reference.NameRegexp.MatchString(mnfst.Name) {
		errs = append(errs,
			distribution.ErrManifestNameInvalid{
				Name:   mnfst.Name,
				Reason: fmt.Errorf("invalid manifest name format"),
			})
	}

	if len(mnfst.History) != len(mnfst.FSLayers) {
		errs = append(errs, fmt.Errorf("mismatched history and fslayer cardinality %d != %d",
			len(mnfst.History), len(mnfst.FSLayers)))
	}

	// TODO: check this
	if _, err := schema1.Verify(mnfst); err != nil {
		switch err {
		case libtrust.ErrMissingSignatureKey, libtrust.ErrInvalidJSONContent, libtrust.ErrMissingSignatureKey:
			errs = append(errs, distribution.ErrManifestUnverified{})
		default:
			if err.Error() == "invalid signature" { // TODO(stevvooe): This should be exported by libtrust
				errs = append(errs, distribution.ErrManifestUnverified{})
			} else {
				errs = append(errs, err)
			}
		}
	}

	// No skipDependencyVerification; always verify
	for _, fsLayer := range mnfst.References() {
		_, err := ms.repository.Blobs(ctx).Stat(ctx, fsLayer.Digest)
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

func (ms *manifestStore) UnmarshalV1(ctx context.Context, dgst digest.Digest, content []byte) (distribution.Manifest, error) {
	context.GetLogger(ms.ctx).Debug("(*signedManifestHandler).Unmarshal")

	var (
		signatures [][]byte
		err        error
	)

	jsig, err := libtrust.NewJSONSignature(content, signatures...)
	if err != nil {
		return nil, err
	}

	// Extract the pretty JWS
	raw, err := jsig.PrettySignature("signatures")
	if err != nil {
		return nil, err
	}

	var sm schema1.SignedManifest
	if err := json.Unmarshal(raw, &sm); err != nil {
		return nil, err
	}
	return &sm, nil
}
