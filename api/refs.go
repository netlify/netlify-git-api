package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
)

// RefUpdateParams is the JSON object sent when patching a ref
type RefUpdateParams struct {
	Sha   string `json:"sha"`
	Force bool   `json:"force"`
}

// GetRef returns a specific reference
func GetRef(w http.ResponseWriter, r *http.Request, params httprouter.Params, ctx context.Context) {
	currentRepo := getRepo(ctx)
	name := path.Join("refs", params.ByName("ref"))
	ref, err := currentRepo.GetRef(name)
	if err != nil {
		HandleError(w, err)
		return
	}

	sendJSON(w, 200, ref)
}

// UpdateRef sets a new target for a reference
func UpdateRef(w http.ResponseWriter, r *http.Request, params httprouter.Params, ctx context.Context) {
	currentRepo := getRepo(ctx)
	refName := path.Join("refs", params.ByName("ref"))
	refParams := &RefUpdateParams{}
	jsonDecoder := json.NewDecoder(r.Body)
	err := jsonDecoder.Decode(refParams)
	if err != nil {
		InternalServerError(w, fmt.Sprintf("Bad parameters to update: %v", err))
		return
	}

	ref, err := currentRepo.UpdateRef(refName, refParams.Sha)
	if err != nil {
		HandleError(w, err)
		return
	}

	sendJSON(w, 200, ref)
}
