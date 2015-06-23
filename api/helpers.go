package api

import (
	"encoding/json"
	"net/http"

	"github.com/netlify/netlify-git-api/repo"
	"golang.org/x/net/context"
)

// Error is an error with a message
type Error struct {
	Msg string `json:"msg"`
}

// InternalServerError sends an error response with a 500 status code
func InternalServerError(w http.ResponseWriter, msg string) {
	sendJSON(w, 500, &Error{Msg: msg})
}

// NotFoundError sends an error response with a 404 status code
func NotFoundError(w http.ResponseWriter, msg string) {
	sendJSON(w, 404, &Error{Msg: msg})
}

// NotAuthorizedError sends an error response with a 401 status code
func NotAuthorizedError(w http.ResponseWriter, msg string) {
	sendJSON(w, 401, &Error{Msg: msg})
}

// HandleError will serve an error response reflecting the error type
func HandleError(w http.ResponseWriter, err error) {
	switch err.(type) {
	default:
		InternalServerError(w, err.Error())
	case *repo.NotFoundError:
		NotFoundError(w, err.Error())
	}
}

func sendJSON(w http.ResponseWriter, status int, obj interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	encoder := json.NewEncoder(w)
	encoder.Encode(obj)
}

// Note - this methods panics if there's no repo in the context
// The context for all handler methods should always have a repo
func getRepo(ctx context.Context) *repo.Repo {
	obj := ctx.Value("repo")
	repo := obj.(*repo.Repo)
	return repo
}
