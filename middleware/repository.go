package middleware

import (
	"github.com/docker/distribution"
	"github.com/docker/distribution/context"
)

// WrappedRepository implements distribution.Repository, providing new calls
// when creating the TagService and MetadataService
type WrappedRepository struct {
	distribution.Repository
}

func (repo *WrappedRepository) Manifests(ctx context.Context, options ...distribution.ManifestServiceOption) (distribution.ManifestService, error) {
	return &manifestStore{
		ctx,
		repo,
		true,
	}, nil
}
