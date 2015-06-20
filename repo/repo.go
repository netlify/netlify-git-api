package repo

import (
	"fmt"

	"github.com/libgit2/git2go"
)

// NotFoundError indicates that an object could not be found
type NotFoundError struct {
	id     string
	object string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("No %v with id %v found", e.object, e.id)
}

// Repo represents the github repo we want to operate on
type Repo struct {
	repo *git.Repository
}

// Open opens a repository
func Open(path string) (*Repo, error) {
	repo, err := git.OpenRepository(path)
	if err != nil {
		return nil, err
	}
	return &Repo{repo}, nil
}
