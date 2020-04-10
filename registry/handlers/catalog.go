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
}
