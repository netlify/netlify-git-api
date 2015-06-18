package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"

	"github.com/julienschmidt/httprouter"
	"github.com/libgit2/git2go"
)

// Reference represents a git reference
type Reference struct {
	Name   string     `json:"ref"`
	Object *RefObject `json:"object"`
}

// RefObject represents the object a reference points to
type RefObject struct {
	Type string `json:"type"`
	Sha  string `json:"sha"`
}

// RefUpdateParams is the JSON object sent when patching a ref
type RefUpdateParams struct {
	Sha   string `json:"sha"`
	Force bool   `json:"force"`
}

// GetRef returns a specific reference
func GetRef(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	refName := path.Join("refs", params.ByName("ref"))
	refObj, err := repo.LookupReference(refName)
	if err != nil {
		NotFoundError(w, fmt.Sprintf("Reference not found: %v", err))
		return
	}

	ref := &Reference{
		Name: refName,
		Object: &RefObject{
			Type: "commit", // This might not always be true?
			Sha:  refObj.Target().String(),
		},
	}

	sendJSON(w, 200, ref)
}

// UpdateRef sets a new target for a reference
func UpdateRef(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	refName := path.Join("refs", params.ByName("ref"))
	refObj, err := repo.LookupReference(refName)
	if err != nil {
		NotFoundError(w, fmt.Sprintf("Reference not found: %v", err))
		return
	}

	refParams := &RefUpdateParams{}
	jsonDecoder := json.NewDecoder(r.Body)
	err = jsonDecoder.Decode(refParams)
	if err != nil {
		InternalServerError(w, fmt.Sprintf("Bad parameters to update: %v", err))
		return
	}

	oid, err := git.NewOid(refParams.Sha)
	if err != nil {
		InternalServerError(w, fmt.Sprintf("Bad sha for update: %v", err))
		return
	}

	refObj, err = refObj.SetTarget(oid, nil, "")
	if err != nil {
		InternalServerError(w, fmt.Sprintf("Could not update ref: %v", err))
		return
	}

	ref := &Reference{
		Name: refName,
		Object: &RefObject{
			Type: "commit", // This might not always be true?
			Sha:  refObj.Target().String(),
		},
	}

	if !repo.IsBare() {
		repo.CheckoutHead(nil)
	}

	sendJSON(w, 200, ref)
}
