package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/libgit2/git2go"
)

// RepoTree a tree in the repository object db
type RepoTree struct {
	Sha  string       `json:"sha"`
	Tree []*TreeEntry `json:"tree,omitempty"`
}

// TreeEntry a single entry in a RepoTree
type TreeEntry struct {
	Path string `json:"path"`
	Type string `json:"type"`
	Mode string `json:"mode"`
	Size int64  `json:"size"`
	Sha  string `json:"sha"`
}

// TreeCreateParams is the JSON object sent when creating a new tree
type TreeCreateParams struct {
	Base string       `json:"base_tree"`
	Tree []*TreeEntry `json:"tree"`
}

func newTreeEntry(repo *git.Repository, entry *git.TreeEntry) *TreeEntry {
	var size int64
	var objType string
	if entry.Type == git.ObjectBlob {
		blob, _ := repo.LookupBlob(entry.Id)
		size = blob.Size()
		objType = "blob"
	} else {
		objType = "tree"
	}
	return &TreeEntry{
		Path: entry.Name,
		Size: size,
		Mode: fmt.Sprintf("%v", entry.Filemode),
		Sha:  entry.Id.String(),
		Type: objType,
	}
}

// CreateTree creates a new tree in the object db
func CreateTree(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	treeParams := &TreeCreateParams{}
	jsonDecoder := json.NewDecoder(r.Body)
	err := jsonDecoder.Decode(treeParams)
	if err != nil {
		InternalServerError(w, fmt.Sprintf("Could not read tree creation params: %v", err))
		return
	}

	builder, err := repo.TreeBuilder()
	if err != nil {
		InternalServerError(w, fmt.Sprintf("Unable to instantiate tree builder: %v", err))
		return
	}

	defer builder.Free()

	if treeParams.Base != "" {
		baseID, err := git.NewOid(treeParams.Base)
		if err != nil {
			InternalServerError(w, fmt.Sprintf("Bad sha for base_tree: %v", err))
			return
		}

		base, err := repo.LookupTree(baseID)
		if err != nil {
			InternalServerError(w, fmt.Sprintf("Could not find any base_tree %v: %v", treeParams.Base, err))
			return
		}

		var i uint64
		for i = 0; i < base.EntryCount(); i++ {
			entry := base.EntryByIndex(i)
			err = builder.Insert(entry.Name, entry.Id, int(entry.Filemode))
			if err != nil {
				InternalServerError(w, fmt.Sprintf("Failed to create tree - base entry %v could not be inserted: %v", entry.Name, err))
				return
			}
		}
	}

	for _, entry := range treeParams.Tree {
		oid, err := git.NewOid(entry.Sha)
		if err != nil {
			InternalServerError(w, fmt.Sprintf("Bad sha for entry %v: %v", entry.Path, err))
			return
		}
		mode, err := strconv.Atoi(entry.Mode)
		if err != nil {
			InternalServerError(w, fmt.Sprintf("Bad file mode for entry %v: %v", entry.Path, err))
			return
		}
		err = builder.Insert(entry.Path, oid, mode)
		if err != nil {
			InternalServerError(w, fmt.Sprintf("Failed to create tree - %v could not be inserted: %v", entry.Path, err))
			return
		}
	}

	oid, err := builder.Write()
	if err != nil {
		InternalServerError(w, fmt.Sprintf("Unable to create tree: %v", err))
		return
	}
	GetTree(w, r, httprouter.Params{httprouter.Param{Key: "sha", Value: oid.String()}})
}

// GetTree gets a representation of a single tree
func GetTree(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	sha := params.ByName("sha")
	oid, err := git.NewOid(sha)
	if err != nil {
		InternalServerError(w, fmt.Sprintf("Invalid sha: %v", err))
		return
	}

	tree, err := repo.LookupTree(oid)
	if err != nil {
		NotFoundError(w, fmt.Sprintf("Tree not found for %v: %v", sha, err))
		return
	}

	repoTree := &RepoTree{
		Sha:  sha,
		Tree: make([]*TreeEntry, tree.EntryCount()),
	}

	var i uint64
	for i = 0; i < tree.EntryCount(); i++ {
		repoTree.Tree[i] = newTreeEntry(repo, tree.EntryByIndex(i))
	}

	sendJSON(w, 200, repoTree)
}
