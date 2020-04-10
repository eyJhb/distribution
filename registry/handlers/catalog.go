package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/handlers"
)

func catalogDispatcher(ctx *Context, r *http.Request) http.Handler {
	catalogHandler := &catalogHandler{
		Context: ctx,
	}

	return handlers.MethodHandler{
		"GET": http.HandlerFunc(catalogHandler.GetCatalog),
	}
}

type catalogHandler struct {
	*Context
}

type catalogAPIResponse struct {
	Repositories []string `json:"repositories"`
}

func (ch *catalogHandler) GetCatalog(w http.ResponseWriter, r *http.Request) {
	ps := paginateStruct{
		errs:      &ch.Errors,
		entriesFn: ch.App.registry.Repositories,
		encFn: func(enc *json.Encoder, n int, repos []string) error {
			return enc.Encode(catalogAPIResponse{
				Repositories: repos[0:n],
			})
		},
	}

	paginateEndpoint(ch.Context, ps, w, r)
	return

	// 	var moreEntries = true

	// 	q := r.URL.Query()
	// 	lastEntry := q.Get("last")
	// 	maxEntries, err := strconv.Atoi(q.Get("n"))
	// 	if err != nil || maxEntries < 0 {
	// 		maxEntries = maximumReturnedEntries
	// 	}

	// 	repos := make([]string, maxEntries)

	// 	filled, err := ch.App.registry.Repositories(ch.Context, repos, lastEntry)
	// 	_, pathNotFound := err.(driver.PathNotFoundError)

	// 	if err == io.EOF || pathNotFound {
	// 		moreEntries = false
	// 	} else if err != nil {
	// 		ch.Errors = append(ch.Errors, errcode.ErrorCodeUnknown.WithDetail(err))
	// 		return
	// 	}
	// 	fmt.Println(filled, err, moreEntries)

	// 	w.Header().Set("Content-Type", "application/json")

	// 	// Add a link header if there are more entries to retrieve
	// 	if moreEntries {
	// 		lastEntry = repos[len(repos)-1]
	// 		urlStr, err := createLinkEntry(r.URL.String(), maxEntries, lastEntry)
	// 		if err != nil {
	// 			ch.Errors = append(ch.Errors, errcode.ErrorCodeUnknown.WithDetail(err))
	// 			return
	// 		}
	// 		w.Header().Set("Link", urlStr)
	// 	}

	// 	enc := json.NewEncoder(w)
	// 	if err := enc.Encode(catalogAPIResponse{
	// 		Repositories: repos[0:filled],
	// 	}); err != nil {
	// 		ch.Errors = append(ch.Errors, errcode.ErrorCodeUnknown.WithDetail(err))
	// 		return
	// 	}
}
