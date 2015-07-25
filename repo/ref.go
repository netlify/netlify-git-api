package repo

import (
	"fmt"
	"log"
	"strings"
	"time"

	"gopkg.in/libgit2/git2go.v22"
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

// GetRef looks up a reference from a name (ie. refs/heads/master)
func (r *Repo) GetRef(name string) (*Reference, error) {
	ref, err := r.repo.LookupReference(name)
	if err != nil {
		return nil, &NotFoundError{id: name, object: "Ref"}
	}

	return &Reference{
		Name: name,
		Object: &RefObject{
			Type: "commit", // This might not always be true?
			Sha:  ref.Target().String(),
		},
	}, nil

}

// UpdateRef updates a reference to point to a new object
// Will check if the repo user has sufficient permissions to
// perform this update
func (r *Repo) UpdateRef(name, newSha string) (*Reference, error) {
	ref, err := r.repo.LookupReference(name)
	if err != nil {
		return nil, &NotFoundError{id: name, object: "Ref"}
	}

	oid, err := git.NewOid(newSha)
	if err != nil {
		return nil, err
	}

	oldCommit, err := r.GetCommit(ref.Target().String())
	if err != nil {
		return nil, err
	}

	newCommit, err := r.GetCommit(newSha)
	if err != nil {
		return nil, err
	}

	changes, err := oldCommit.ChangedFiles(newCommit)
	if err != nil {
		return nil, err
	}

	failMsg := []string{}
	for _, change := range changes {
		if !r.user.HasPermission(change.Action, change.Path) {
			failMsg = append(failMsg, fmt.Sprintf("you do not have permission to %v: %v", change.Action, change.Path))
		}
	}

	if len(failMsg) > 0 {
		return nil, &ForbiddenError{msg: strings.Join(failMsg, ",")}
	}

	sig := &git.Signature{
		Name:  r.user.Name(),
		Email: r.user.Email(),
		When:  time.Now(),
	}

	ref, err = ref.SetTarget(oid, sig, "")
	if err != nil {
		return nil, err
	}

	if !r.repo.IsBare() {
		paths := make([]string, len(changes))
		del := false
		for i, change := range changes {
			paths[i] = change.Path
			if change.Action == "delete" {
				del = true
			}
		}

		tree, _ := r.repo.LookupTree(newCommit.Tree.id)
		options := &git.CheckoutOpts{Strategy: git.CheckoutForce}
		// If we're deleting files, passing paths to options makes those deleted files
		// stick around in the index and causes weirdness, so we only specify paths (safer)
		// when there are no deletes
		if del == false {
			options.Paths = paths
		}
		err = r.repo.CheckoutTree(tree, options)
		if err != nil {
			log.Fatalf("Error checking out head: %v", err)
		}
	}

	return &Reference{
		Name: name,
		Object: &RefObject{
			Type: "commit", // This might not always be true?
			Sha:  ref.Target().String(),
		},
	}, nil
}
