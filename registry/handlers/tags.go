package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/handlers"
)

// tagsDispatcher constructs the tags handler api endpoint.
func tagsDispatcher(ctx *Context, r *http.Request) http.Handler {
	tagsHandler := &tagsHandler{
		Context: ctx,
	}

	return handlers.MethodHandler{
		"GET": http.HandlerFunc(tagsHandler.GetTags),
	}
}

// tagsHandler handles requests for lists of tags under a repository name.
type tagsHandler struct {
	*Context
}

type tagsAPIResponse struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

// GetTags returns a json list of tags for a specific image name.
func (th *tagsHandler) GetTags(w http.ResponseWriter, r *http.Request) {
	tagService := th.Repository.Tags(th)

	ps := paginateStruct{
		errs:      &th.Errors,
		entriesFn: tagService.Tags,
		encFn: func(enc *json.Encoder, n int, tags []string) error {
			return enc.Encode(tagsAPIResponse{
				Name: th.Repository.Named().Name(),
				Tags: tags[0:n],
			})
		},
	}

	paginateEndpoint(th.Context, ps, w, r)
	return
}
