package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/netlify/netlify-git-api/repo"
	"golang.org/x/net/context"
)

// FileDeleteParams holds the parameters for DeleteFile request
type FileDeleteParams struct {
	Sha     string `json:"sha"`
	Message string `json:"message"`
	Branch  string `json:"branch"`
}

// GetFile returns information about a file or directory in the repository.
// If the Content-Type is set to "application/vnd.netlify.raw" it will return
// the actual file contents (or an error if a directory)
func GetFile(w http.ResponseWriter, r *http.Request, params httprouter.Params, ctx context.Context) {
	currentRepo := getRepo(ctx)
	pathname := params.ByName("path")[1:]
	file, err := currentRepo.GetFile(pathname)
	if err != nil {
		HandleError(w, err)
		return
	}

	if file.Type == "dir" {
		sendJSON(w, 200, file.Files)
		return
	}

	if r.Header.Get("Content-Type") != rawContentType {
		sendJSON(w, 200, file)
	}

	blob, err := currentRepo.GetBlob(file.Sha)
	if err != nil {
		HandleError(w, err)
		return
	}
	io.Copy(w, blob)
}

// DeleteFile deletes a file from the repo
// Takes a `sha`, a `message` for the commit message and the path to the file
func DeleteFile(w http.ResponseWriter, r *http.Request, params httprouter.Params, ctx context.Context) {
	currentRepo := getRepo(ctx)
	pathname := params.ByName("path")[1:]

	fileParams := &FileDeleteParams{}
	jsonDecoder := json.NewDecoder(r.Body)
	err := jsonDecoder.Decode(fileParams)
	if err != nil {
		InternalServerError(w, fmt.Sprintf("Could not read file deletiong params: %v", err))
		return
	}

	ref, err := currentRepo.GetRef("refs/heads/" + fileParams.Branch)
	if err != nil {
		HandleError(w, err)
		return
	}
	commit, err := currentRepo.GetCommit(ref.Object.Sha)
	if err != nil {
		HandleError(w, err)
		return
	}

	tree, err := currentRepo.GetTree(commit.Tree.Sha)
	if err != nil {
		HandleError(w, err)
		return
	}

	dir, name := path.Split(pathname)
	segments := strings.Split(dir, "/")
	trees := []*repo.Tree{tree}
	for _, segment := range segments {
		if segment == "" {
			continue
		}
		for _, entry := range tree.Tree {
			if entry.Path == segment {
				tree, err = currentRepo.GetTree(entry.Sha)
				if err != nil {
					HandleError(w, err)
					return
				}
				trees = append(trees, tree)
				break
			}
		}
	}

	newEntries := []*repo.TreeEntry{}
	for _, entry := range tree.Tree {
		if entry.Path != name {
			newEntries = append(newEntries, entry)
		}
	}
	newTree, err := currentRepo.CreateTree("", newEntries)

	for i := len(trees) - 1; i > 0; i-- {
		tree = trees[i-1]
		name = segments[i-1]
		for _, entry := range tree.Tree {
			if entry.Path == name {
				entry.Sha = newTree.Sha
			}
		}
		newTree, err = currentRepo.CreateTree("", tree.Tree)
	}

	newCommit, err := currentRepo.CreateCommit(newTree.Sha, fileParams.Message, []string{ref.Object.Sha})
	if err != nil {
		HandleError(w, err)
		return
	}

	newRef, err := currentRepo.UpdateRef("refs/heads/"+fileParams.Branch, newCommit.Sha)

	if err != nil {
		HandleError(w, err)
		return
	}
	sendJSON(w, 200, newRef)
}
