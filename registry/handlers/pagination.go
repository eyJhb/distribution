package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/docker/distribution/registry/api/errcode"
	"github.com/docker/distribution/registry/storage/driver"
)

const maximumReturnedEntries = 100

type paginateStruct struct {
	errs      *errcode.Errors
	entriesFn func(ctx context.Context, results []string, lastEntry string) (n int, err error)
	encFn     func(enc *json.Encoder, n int, results []string) error
}

func paginateEndpoint(ctx context.Context, ps paginateStruct, w http.ResponseWriter, r *http.Request) {
	var moreEntries = true

	q := r.URL.Query()
	lastEntry := q.Get("last")
	maxEntries, err := strconv.Atoi(q.Get("n"))
	if err != nil || maxEntries < 0 {
		maxEntries = maximumReturnedEntries
	}

	results := make([]string, maxEntries)

	n, err := ps.entriesFn(ctx, results, lastEntry)
	_, pathNotFound := err.(driver.PathNotFoundError)

	if err == io.EOF || pathNotFound {
		moreEntries = false
	} else if err != nil {
		*ps.errs = append(*ps.errs, errcode.ErrorCodeUnknown.WithDetail(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Add a link header if there are more entries to retrieve
	if moreEntries {
		lastEntry = results[len(results)-1]
		urlStr, err := createLinkEntry(r.URL.String(), maxEntries, lastEntry)
		if err != nil {
			*ps.errs = append(*ps.errs, errcode.ErrorCodeUnknown.WithDetail(err))
			return
		}
		w.Header().Set("Link", urlStr)
	}

	enc := json.NewEncoder(w)
	if err := ps.encFn(enc, n, results); err != nil {
		*ps.errs = append(*ps.errs, errcode.ErrorCodeUnknown.WithDetail(err))
	}

	return
}

// Use the original URL from the request to create a new URL for
// the link header
func createLinkEntry(origURL string, maxEntries int, lastEntry string) (string, error) {
	calledURL, err := url.Parse(origURL)
	if err != nil {
		return "", err
	}

	v := url.Values{}
	v.Add("n", strconv.Itoa(maxEntries))
	v.Add("last", lastEntry)

	calledURL.RawQuery = v.Encode()

	calledURL.Fragment = ""
	urlStr := fmt.Sprintf("<%s>; rel=\"next\"", calledURL.String())

	return urlStr, nil
}
