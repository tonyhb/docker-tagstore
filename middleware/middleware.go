package middleware

import (
	"github.com/docker/distribution"
	"github.com/docker/distribution/context"
)

func InitMiddleware(ctx context.Context, repository distribution.Repository, options map[string]interface{}) (distribution.Repository, error) {

	// TODO: Expose all registry config items in here. They're necessary
	// for your middleware to get things right.
	//
	// In short, these are the things we need for manifest services:
	//
	//
	// TODO: expose whether signatures are enabled (used in PUT and GET)
	// TODO: expose schema1signingkey
	// TODO: expose whether delete is enabled within the middleware here.
	// TODO: expose whether this is a pull-through cache in the middleware
	// options; the ManifestService has a `skipDependencyVerification`
	// setting which verifies that layers exist in the blobstore when
	// saving manifests.
	//
	// BUT - is this necessary? a pull-through cache should only pull...

	return &WrappedRepository{
		Repository: repository,
	}, nil

}
