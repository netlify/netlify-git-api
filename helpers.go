package main

import (
	"encoding/json"
	"net/http"
)

// InternalServerError sends an error response with a 500 status code
func InternalServerError(w http.ResponseWriter, msg string) {
	sendJSON(w, 500, &Error{Msg: msg})
}

// NotFoundError sends an error response with a 404 status code
func NotFoundError(w http.ResponseWriter, msg string) {
	sendJSON(w, 404, &Error{Msg: msg})
}

func sendJSON(w http.ResponseWriter, status int, obj interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	encoder := json.NewEncoder(w)
	encoder.Encode(obj)
}
