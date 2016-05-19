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

// signedManifestHandler is a ManifestHandler that covers schema1 manifests. It
// can unmarshal and put schema1 manifests that have been signed by libtrust.
type signedManifestHandler struct {
	repository distribution.Repository
	ctx        context.Context
	signatures *signatureStore
}

// verifyManifest ensures that the manifest content is valid from the
// perspective of the registry. It ensures that the signature is valid for the
// enclosed payload. As a policy, the registry only tries to store valid
// content, leaving trust policies of that content up to consumers.
func (ms *signedManifestHandler) verifyManifest(ctx context.Context, mnfst schema1.SignedManifest, skipDependencyVerification bool) error {
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

	if _, err := schema1.Verify(&mnfst); err != nil {
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

	if !skipDependencyVerification {
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
	}
	if len(errs) != 0 {
		return errs
	}

	return nil
}
