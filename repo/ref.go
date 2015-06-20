package repo

import "github.com/libgit2/git2go"

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
func (r *Repo) UpdateRef(name, newSha string) (*Reference, error) {
	ref, err := r.repo.LookupReference(name)
	if err != nil {
		return nil, &NotFoundError{id: name, object: "Ref"}
	}

	oid, err := git.NewOid(newSha)
	if err != nil {
		return nil, err
	}

	ref, err = ref.SetTarget(oid, nil, "")
	if err != nil {
		return nil, err
	}

	return &Reference{
		Name: name,
		Object: &RefObject{
			Type: "commit", // This might not always be true?
			Sha:  ref.Target().String(),
		},
	}, nil
}
