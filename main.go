package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
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
	Size int    `json:"size"`
	Sha  string `json:"sha"`
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
)

func init() {
	commits = map[string]Commit{}
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

	fmt.Fprint(w, "Not implemented")
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

	log.Fatal(http.ListenAndServe(":8080", router))
}
