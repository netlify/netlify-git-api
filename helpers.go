package main

import (
	"encoding/json"
	"net/http"

	gitrepo "github.com/netlify/netlify-git-api/repo"
)

// InternalServerError sends an error response with a 500 status code
func InternalServerError(w http.ResponseWriter, msg string) {
	sendJSON(w, 500, &Error{Msg: msg})
}

// NotFoundError sends an error response with a 404 status code
func NotFoundError(w http.ResponseWriter, msg string) {
	sendJSON(w, 404, &Error{Msg: msg})
}

// HandleError will serve an error response reflecting the error type
func HandleError(w http.ResponseWriter, err error) {
	switch err.(type) {
	default:
		InternalServerError(w, err.Error())
	case *gitrepo.NotFoundError:
		NotFoundError(w, err.Error())
	}
}

func sendJSON(w http.ResponseWriter, status int, obj interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	encoder := json.NewEncoder(w)
	encoder.Encode(obj)
}
