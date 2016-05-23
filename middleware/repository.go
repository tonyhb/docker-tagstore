package middleware

import (
	"github.com/docker/distribution"
	"github.com/docker/distribution/context"
)

// WrappedRepository implements distribution.Repository, providing new calls
// when creating the TagService and MetadataService
type WrappedRepository struct {
	distribution.Repository

	store Store
}

func (repo *WrappedRepository) Manifests(ctx context.Context, options ...distribution.ManifestServiceOption) (distribution.ManifestService, error) {
	return &manifestStore{
		ctx,
		repo,
		repo.store,
	}, nil
}

func (repo *WrappedRepository) Tags(ctx context.Context) distribution.TagService {
	return &tagStore{
		ctx,
		repo,
		repo.store,
	}
}
