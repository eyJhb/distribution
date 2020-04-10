package storage

import (
	"context"
	"path"

	"github.com/docker/distribution/reference"
	"github.com/docker/distribution/registry/storage/driver"
)

// Returns a list, or partial list, of repositories in the registry.
// Because it's a quite expensive operation, it should only be used when building up
// an initial set of repositories.
func (reg *registry) Repositories(ctx context.Context, repos []string, last string) (n int, err error) {
	root, err := pathFor(repositoriesRootPathSpec{})
	if err != nil {
		return 0, err
	}

	return paginateEndpoint(ctx, repos, root, "_layers", last, reg.blobStore.driver.Walk)
}

// Enumerate applies ingester to each repository
func (reg *registry) Enumerate(ctx context.Context, ingester func(string) error) error {
	root, err := pathFor(repositoriesRootPathSpec{})
	if err != nil {
		return err
	}

	err = reg.blobStore.driver.Walk(ctx, root, func(fileInfo driver.FileInfo) error {
		return handlePaginateWalk(fileInfo, root, "_layers", "", ingester)
	})

	return err
}

// Remove removes a repository from storage
func (reg *registry) Remove(ctx context.Context, name reference.Named) error {
	root, err := pathFor(repositoriesRootPathSpec{})
	if err != nil {
		return err
	}
	repoDir := path.Join(root, name.Name())
	return reg.driver.Delete(ctx, repoDir)
}
