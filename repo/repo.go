package repo

import (
	"fmt"

	"gopkg.in/libgit2/git2go.v22"
)

// NotFoundError indicates that an object could not be found
type NotFoundError struct {
	id     string
	object string
}

// ForbiddenError indicates that the user doesn't have permission to do this action
type ForbiddenError struct {
	msg string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("No %v with id %v found", e.object, e.id)
}

func (e *ForbiddenError) Error() string {
	return e.msg
}

// Repo represents the github repo we want to operate on
type Repo struct {
	repo *git.Repository
	user User
	sync bool
}

// User is the main user object for the API.
type User interface {
	Name() string
	Email() string
	HasPermission(string, string) bool
}

// Open opens a repository
func Open(user User, path string, sync bool) (*Repo, error) {
	repo, err := git.OpenRepository(path)
	if err != nil {
		return nil, err
	}

	return &Repo{repo: repo, user: user, sync: sync}, nil
}
