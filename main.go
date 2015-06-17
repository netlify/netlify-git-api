package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/julienschmidt/httprouter"
	git "github.com/libgit2/git2go"
	"github.com/rs/cors"
)

// CommitFile represents a single file to be created, updated or removed when a commit
// is finalized
type CommitFile struct {
	Path    string
	Content []byte
}

// Commit is an open commit. When applied the AddedFiles will be updated or created
// in the current directory, the RemovedFiles will be removed from the directory
// and the commit will then be committed with the specified commit Message
type Commit struct {
	ID           string
	Message      string
	AddedFiles   []CommitFile
	RemovedFiles []CommitFile
}

// RepoFile a file in the repository
type RepoFile struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Size int64  `json:"size"`
	Sha  string `json:"sha"`
	Type string `json:"type"`
}

// Error is an error with a message
type Error struct {
	Msg string `json:"msg"`
}

var (
	indexPage = `<!doctype html>
<html>
  <head><title>Local Netlify CMS Backend</title></head>
  <body>
    <h1>Local Netlify CMS Backend</h1>
    <p>
      This is a simple local backend for <a href="https://github.com/netlify/cms">Netlify's CMS</a>
    </p>
  </body>
</html>
  `
	commits map[string]Commit
	repo    *git.Repository
)

func init() {
	commits = map[string]Commit{}
	cwd, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("Error getting current working dir: %v", err))
	}
	repo, err = git.OpenRepository(cwd)
	if err != nil {
		panic(fmt.Sprintf("Unable to open git repository in %v: %v", cwd, err))
	}
}

func sendJSON(w http.ResponseWriter, status int, obj interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	encoder := json.NewEncoder(w)
	encoder.Encode(obj)
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

// InternalServerError sends an error response with a 500 status code
func InternalServerError(w http.ResponseWriter, msg string) {
	sendJSON(w, 500, &Error{Msg: msg})
}

// NotFoundError sends an error response with a 404 status code
func NotFoundError(w http.ResponseWriter, msg string) {
	sendJSON(w, 404, &Error{Msg: msg})
}

// Index serves a basic index page
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, indexPage)
}

// Auth should be opened in a popup. Will create a new token and return it to
// the origin window via postMessage
func Auth(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Not implemented")
}

// GetFile returns information about a file or directory in the repository.
// If the Content-Type is set to "application/vnd.netlify.raw" it will return
// the actual file contents (or an error if a directory)
func GetFile(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	objPath := params.ByName("path")
	objPath = objPath[1:]

	ref, err := repo.Head()
	if err != nil {
		InternalServerError(w, fmt.Sprintf("Unable to get head for repo: %v", err))
		return
	}

	obj, err := repo.Lookup(ref.Target())
	if err != nil {
		InternalServerError(w, fmt.Sprintf("Unable to get tree for head: %v", err))
		return
	}

	commit, ok := obj.(*git.Commit)
	if !ok {
		InternalServerError(w, fmt.Sprintf("Head is not a commit: %v", obj.Type()))
		return
	}

	tree, err := commit.Tree()
	if err != nil {
		InternalServerError(w, fmt.Sprintf("Error getting tree for commit: %v", err))
		return
	}

	log.Printf("Looking up path %v in tree %v", objPath, tree)
	var entry *git.TreeEntry
	if objPath == "" {
		entry = &git.TreeEntry{Name: "", Id: tree.Id(), Type: git.ObjectTree}
	} else {
		entry, err = tree.EntryByPath(objPath)
	}

	if err != nil {
		NotFoundError(w, fmt.Sprintf("Unable to find path %v. %v", objPath, err))
		return
	}

	if entry.Type == git.ObjectTree {
		tree, err = repo.LookupTree(entry.Id)
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

		if r.Header.Get("Content-Type") == "application/vnd.netlify.raw" {
			w.Write(blob.Contents())
			return
		}

		file := newRepoFile(repo, entry, path.Dir(objPath))
		sendJSON(w, 200, file)
	}
}

// CreateCommit creates a new open commit and returns a commit ID
func CreateCommit(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Not implemented")
}

// AddFile ads a file to an open commit
func AddFile(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Not implemented")
}

// RemoveFile removes a file from an open commit
func RemoveFile(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Not implemented")
}

// FinalizeCommit applies a commit to the repository, takes a `message` parameter
// that will set the commit message.
func FinalizeCommit(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Not implemented")
}

func main() {
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/auth", Auth)
	router.GET("/files/*path", GetFile)
	router.POST("/commits", CreateCommit)
	router.PUT("/commits/:id/*path", AddFile)
	router.DELETE("/commits/:id/*path", RemoveFile)
	router.POST("/commits/:id", FinalizeCommit)

	handler := cors.Default().Handler(router)

	log.Fatal(http.ListenAndServe(":8080", handler))
}
