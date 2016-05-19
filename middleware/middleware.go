package middleware

import (
	"github.com/docker/distribution"
	"github.com/docker/distribution/context"
)

func InitMiddleware(ctx context.Context, repository distribution.Repository, options map[string]interface{}) (distribution.Repository, error) {
	// TODO: expose whether delete is enabled within the middleware here.
	// if this is not an option we must create a RegistryMiddleware which has
	// access to distribution.Namespace

	return &WrappedRepository{
		Repository: repository,
	}, nil

}
