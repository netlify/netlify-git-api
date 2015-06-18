package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/libgit2/git2go"
)

// Commit represents a commit
type Commit struct {
	Sha       string          `json:"sha"`
	Author    *Author         `json:"author,omitempty"`
	Committer *Author         `json:"committer,omitempty"`
	Message   string          `json:"message"`
	Tree      *RepoTree       `json:"tree"`
	Parents   []*CommitParent `json:"parents"`
}

// CommitParent represents a parent of a commit
type CommitParent struct {
	Sha string `json:"sha"`
}

// Author of a commit
type Author struct {
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Date  time.Time `json:"date"`
}

// CommitCreateParams is the JSON object sent when creating a new commit
type CommitCreateParams struct {
	Msg     string   `json:"message"`
	Tree    string   `json:"tree"`
	Parents []string `json:"parents"`
}

// CreateCommit creates a new commit
func CreateCommit(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	commitParams := &CommitCreateParams{}
	jsonDecoder := json.NewDecoder(r.Body)
	err := jsonDecoder.Decode(commitParams)
	if err != nil {
		InternalServerError(w, fmt.Sprintf("Could not read commit creation params: %v", err))
		return
	}

	treeID, err := git.NewOid(commitParams.Tree)
	if err != nil {
		InternalServerError(w, fmt.Sprintf("Bad sha for commit tree: %v", err))
		return
	}

	tree, err := repo.LookupTree(treeID)
	if err != nil {
		InternalServerError(w, fmt.Sprintf("Could not find any tree %v: %v", treeID.String(), err))
		return
	}

	parents := make([]*git.Commit, len(commitParams.Parents))
	for i := 0; i < len(parents); i++ {
		oid, err := git.NewOid(commitParams.Parents[i])
		if err != nil {
			InternalServerError(w, fmt.Sprintf("Bad commit id %v: %v", oid, err))
			return
		}

		commit, err := repo.LookupCommit(oid)
		if err != nil {
			InternalServerError(w, fmt.Sprintf("Could not find any commit %v: %v", oid.String(), err))
			return
		}
		parents[i] = commit
	}

	sig := &git.Signature{
		Name:  name,
		Email: email,
		When:  time.Now(),
	}

	oid, err := repo.CreateCommit("", sig, sig, commitParams.Msg, tree, parents...)
	if err != nil {
		InternalServerError(w, fmt.Sprintf("Error creating commit: %v", err))
		return
	}

	GetCommit(w, r, httprouter.Params{httprouter.Param{Key: "sha", Value: oid.String()}})
}

// GetCommit returns a single commit object
func GetCommit(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	sha := params.ByName("sha")
	oid, err := git.NewOid(sha)
	if err != nil {
		InternalServerError(w, fmt.Sprintf("Invalid sha: %v", err))
		return
	}

	commit, err := repo.LookupCommit(oid)
	if err != nil {
		NotFoundError(w, fmt.Sprintf("Commit not found for %v: %v", sha, err))
		return
	}

	repoCommit := &Commit{
		Sha:     sha,
		Message: commit.Message(),
		Tree:    &RepoTree{Sha: commit.TreeId().String()},
		Parents: make([]*CommitParent, commit.ParentCount()),
	}

	author := commit.Author()
	if author != nil {
		repoCommit.Author = &Author{Name: author.Name, Email: author.Email, Date: author.When}
	}

	committer := commit.Committer()
	if committer != nil {
		repoCommit.Committer = &Author{Name: committer.Name, Email: committer.Email, Date: committer.When}
	}

	var i uint
	for i = 0; i < commit.ParentCount(); i++ {
		repoCommit.Parents[i] = &CommitParent{Sha: commit.ParentId(i).String()}
	}

	sendJSON(w, 200, repoCommit)
}
