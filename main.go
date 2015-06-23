package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/netlify/netlify-git-api/api"
	"github.com/netlify/netlify-git-api/repo"
)

type user struct {
	name  string
	email string
}

func (u *user) Name() string {
	return u.name
}

func (u *user) Email() string {
	return u.email
}

func (u *user) HasPermission(_ string, _ string) bool {
	return true
}

type resolver struct {
	repo *repo.Repo
}

func (r *resolver) GetRepo(_ *http.Request) (*repo.Repo, error) {
	return r.repo, nil
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("Error getting current working dir: %v", err))
	}

	user := &user{name: "Netlify CMS", email: "team@netlify.com"}

	currentRepo, err := repo.Open(user, cwd)
	if err != nil {
		panic(fmt.Sprintf("Unable to open git repository in %v: %v", cwd, err))
	}

	resolver := &resolver{repo: currentRepo}

	api := api.NewAPI(resolver)
	log.Fatal(http.ListenAndServe(":8080", api))
}
