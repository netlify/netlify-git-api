package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	gitrepo "github.com/netlify/netlify-git-api/repo"
	"golang.org/x/net/context"
)

// TreeCreateParams is the JSON object sent when creating a new tree
type TreeCreateParams struct {
	Base string               `json:"base_tree"`
	Tree []*gitrepo.TreeEntry `json:"tree"`
}

// CreateTree creates a new tree in the object db
func CreateTree(w http.ResponseWriter, r *http.Request, params httprouter.Params, ctx context.Context) {
	currentRepo := getRepo(ctx)
	treeParams := &TreeCreateParams{}
	jsonDecoder := json.NewDecoder(r.Body)
	err := jsonDecoder.Decode(treeParams)
	if err != nil {
		InternalServerError(w, fmt.Sprintf("Could not read tree creation params: %v", err))
		return
	}

	tree, err := currentRepo.CreateTree(treeParams.Base, treeParams.Tree)
	if err != nil {
		HandleError(w, err)
	}

	sendJSON(w, 200, tree)
}

// GetTree gets a representation of a single tree
func GetTree(w http.ResponseWriter, r *http.Request, params httprouter.Params, ctx context.Context) {
	currentRepo := getRepo(ctx)
	tree, err := currentRepo.GetTree(params.ByName("sha"))
	if err != nil {
		HandleError(w, err)
		return
	}

	sendJSON(w, 200, tree)
}
