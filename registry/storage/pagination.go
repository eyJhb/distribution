package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path"

	"github.com/docker/distribution/registry/storage/driver"
)

var errPaginateStop = errors.New("we have reached last, skip the rest")

// Returns a list, or partial list, of the desired results from drverWalkFn.
func paginateEndpoint(ctx context.Context, results []string, root, resultsKey, last string, drverWalkFn func(ctx context.Context, path string, f driver.WalkFn) error) (n int, err error) {
	var moreResults bool
	var foundResults []string

	if len(results) == 0 {
		return 0, errors.New("no space in slice")
	}

	err = drverWalkFn(ctx, root, func(fileInfo driver.FileInfo) error {
		err := handlePaginateWalk(fileInfo, root, resultsKey, last, func(resultPath string) error {
			// if we have already filed up the `foundResults`, then
			// this means that there are more results
			if len(foundResults) == len(results) {
				moreResults = true
				return errPaginateStop
			}

			foundResults = append(foundResults, resultPath)
			return nil
		})
		if err != nil {
			return err
		}

		return nil
	})

	n = copy(results, foundResults)

	// if moreResults is `true`, then the error
	// is errPaginateStop, which we do not care
	// about here
	if err != nil && !moreResults {
		return n, err
	} else if !moreResults {
		// No more records are available.
		return n, io.EOF
	}

	return n, nil
}

// handlePaginateWalk calls function fn with a repository path if fileInfo
// has a path of a repository under root and that it is lexographically
// after last. Otherwise, it will return ErrSkipDir. This should be used
// with Walk to do handling with repositories in a storage.
func handlePaginateWalk(fileInfo driver.FileInfo, root, resultsKey, last string, fn func(repoPath string) error) error {
	rootFilePath := fileInfo.Path()
	fmt.Println(rootFilePath)

	// lop the base path off
	filePath := rootFilePath[len(root)+1:]

	_, file := path.Split(filePath)
	if last < file {
		if err := fn(file); err != nil {
			return err
		}
	}
	return driver.ErrSkipDir
}
