package api

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Auth should be opened in a popup. Will create a new token and return it to
// the origin window via postMessage
func Auth(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Not implemented")
}
