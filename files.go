package main

import (
	"fmt"
	"log"
	"net/http"
	"path"

	"github.com/julienschmidt/httprouter"
	"github.com/libgit2/git2go"
)

// RepoFile a file in the repository
type RepoFile struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Size int64  `json:"size"`
	Sha  string `json:"sha"`
	Type string `json:"type"`
}

// GetFile returns information about a file or directory in the repository.
// If the Content-Type is set to "application/vnd.netlify.raw" it will return
// the actual file contents (or an error if a directory)
func GetFile(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	objPath := params.ByName("path")[1:]

	entry := findEntry(w, objPath)
	if entry == nil {
		return
	}

	if entry.Type == git.ObjectTree {
		tree, err := repo.LookupTree(entry.Id)
		if err != nil {
			InternalServerError(w, fmt.Sprintf("Error while getting tree: %v", err))
			return
		}

		files := make([]*RepoFile, tree.EntryCount())
		i := 0
		tree.Walk(func(name string, entry *git.TreeEntry) int {
			files[i] = newRepoFile(repo, entry, objPath)
			i++
			return i
		})

		sendJSON(w, 200, files)
	} else {
		blob, err := repo.LookupBlob(entry.Id)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "Error while getting blob: %v", err)
			return
		}
		log.Printf("Type: %v", r.Header.Get("Content-Type"))

		if r.Header.Get("Content-Type") == rawContentType {
			w.Write(blob.Contents())
			return
		}

		file := newRepoFile(repo, entry, path.Dir(objPath))
		sendJSON(w, 200, file)
	}
}

func findEntry(w http.ResponseWriter, objPath string) (entry *git.TreeEntry) {
	ref, err := repo.Head()
	if err != nil {
		InternalServerError(w, fmt.Sprintf("Unable to get head for repo: %v", err))
		return nil
	}

	obj, err := repo.Lookup(ref.Target())
	if err != nil {
		InternalServerError(w, fmt.Sprintf("Unable to get tree for head: %v", err))
		return nil
	}

	commit, ok := obj.(*git.Commit)
	if !ok {
		InternalServerError(w, fmt.Sprintf("Head is not a commit: %v", obj.Type()))
		return nil
	}

	tree, err := commit.Tree()
	if err != nil {
		InternalServerError(w, fmt.Sprintf("Error getting tree for commit: %v", err))
		return nil
	}

	if objPath == "" {
		entry = &git.TreeEntry{Name: "", Id: tree.Id(), Type: git.ObjectTree}
	} else {
		entry, err = tree.EntryByPath(objPath)
	}

	if err != nil {
		NotFoundError(w, fmt.Sprintf("Unable to find path %v. %v", objPath, err))
		return nil
	}

	return entry
}

func newRepoFile(repo *git.Repository, entry *git.TreeEntry, dir string) *RepoFile {
	var size int64
	var objType string
	if entry.Type == git.ObjectBlob {
		blob, _ := repo.LookupBlob(entry.Id)
		size = blob.Size()
		objType = "file"
	} else {
		objType = "dir"
	}
	return &RepoFile{
		Name: entry.Name,
		Path: path.Join(dir, entry.Name),
		Size: size,
		Sha:  entry.Id.String(),
		Type: objType,
	}
}
