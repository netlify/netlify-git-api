package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/netlify/netlify-git-api/repo"
	"github.com/rs/cors"
	"golang.org/x/net/context"
)

const (
	rawContentType = "application/vnd.netlify.raw"
)

// API is the REST API around git repos
type API struct {
	resolver Resolver
}

// Resolver handlers user and repo lookups for requests
type Resolver interface {
	GetRepo(*http.Request) (*repo.Repo, error)
}

func (a *API) wrap(fn func(http.ResponseWriter, *http.Request, httprouter.Params, context.Context)) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		repo, err := a.resolver.GetRepo(r)
		if err != nil {
			HandleError(w, err)
			return
		}

		if repo == nil {
			NotAuthorizedError(w, "No repository resolved")
			return
		}

		ctx := context.WithValue(nil, "repo", repo)

		fn(w, r, p, ctx)
	}
}

// NewAPI instantiates a new API with a resolver
func NewAPI(resolver Resolver) http.Handler {
	api := API{resolver: resolver}
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/auth", Auth)
	router.GET("/files/*path", api.wrap(GetFile))

	router.POST("/blobs", api.wrap(CreateBlob))
	router.GET("/blobs/:sha", api.wrap(GetBlob))

	router.POST("/trees", api.wrap(CreateTree))
	router.GET("/trees/:sha", api.wrap(GetTree))

	router.POST("/commits", api.wrap(CreateCommit))
	router.GET("/commits/:sha", api.wrap(GetCommit))

	router.GET("/refs/*ref", api.wrap(GetRef))
	router.PATCH("/refs/*ref", api.wrap(UpdateRef))

	corsHandler := cors.New(cors.Options{AllowedMethods: []string{"GET", "POST", "PATCH", "PUT", "DELETE"}})

	return corsHandler.Handler(router)
}
