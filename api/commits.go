package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
)

// CommitCreateParams is the JSON object sent when creating a new commit
type CommitCreateParams struct {
	Msg     string   `json:"message"`
	Tree    string   `json:"tree"`
	Parents []string `json:"parents"`
}

// CreateCommit creates a new commit
func CreateCommit(w http.ResponseWriter, r *http.Request, params httprouter.Params, ctx context.Context) {
	currentRepo := getRepo(ctx)
	commitParams := &CommitCreateParams{}
	jsonDecoder := json.NewDecoder(r.Body)
	err := jsonDecoder.Decode(commitParams)
	if err != nil {
		InternalServerError(w, fmt.Sprintf("Could not read commit creation params: %v", err))
		return
	}

	commit, err := currentRepo.CreateCommit(
		commitParams.Tree,
		commitParams.Msg,
		commitParams.Parents,
	)
	if err != nil {
		HandleError(w, err)
		return
	}

	sendJSON(w, 200, commit)
}

// GetCommit returns a single commit object
func GetCommit(w http.ResponseWriter, r *http.Request, params httprouter.Params, ctx context.Context) {
	currentRepo := getRepo(ctx)
	sha := params.ByName("sha")
	commit, err := currentRepo.GetCommit(sha)
	if err != nil {
		HandleError(w, err)
		return
	}

	sendJSON(w, 200, commit)
}
