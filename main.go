package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	git "github.com/libgit2/git2go"
	"github.com/rs/cors"
)

const (
	rawContentType = "application/vnd.netlify.raw"
	name           = "Netlify CMS"
	email          = "team@netlify.com"
)

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
	repo *git.Repository
)

func init() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("Error getting current working dir: %v", err))
	}
	repo, err = git.OpenRepository(cwd)
	if err != nil {
		panic(fmt.Sprintf("Unable to open git repository in %v: %v", cwd, err))
	}
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

func main() {
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/auth", Auth)
	router.GET("/files/*path", GetFile)

	router.POST("/blobs", CreateBlob)
	router.GET("/blobs/:sha", GetBlob)

	router.POST("/trees", CreateTree)
	router.GET("/trees/:sha", GetTree)

	router.POST("/commits", CreateCommit)
	router.GET("/commits/:sha", GetCommit)

	router.GET("/refs/*ref", GetRef)
	router.PATCH("/refs/*ref", UpdateRef)

	corsHandler := cors.New(cors.Options{AllowedMethods: []string{"GET", "POST", "PATCH", "PUT", "DELETE"}})

	log.Fatal(http.ListenAndServe(":8080", corsHandler.Handler(router)))
}
