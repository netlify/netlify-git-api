package repo

import (
	"time"

	"gopkg.in/libgit2/git2go.v22"
)

// Commit represents a commit
type Commit struct {
	id        *git.Oid
	Sha       string          `json:"sha"`
	Author    *Author         `json:"author,omitempty"`
	Committer *Author         `json:"committer,omitempty"`
	Message   string          `json:"message"`
	Tree      *Tree           `json:"tree"`
	Parents   []*CommitParent `json:"parents"`
	repo      *Repo
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

// FileChange represents a file that will change between two commits
// Action can be "create", "update", "delete"
type FileChange struct {
	Action string
	Path   string
}

// GetCommit looks up a commit from a sha
func (r *Repo) GetCommit(sha string) (*Commit, error) {
	oid, err := git.NewOid(sha)
	if err != nil {
		return nil, err
	}

	commit, err := r.repo.LookupCommit(oid)
	if err != nil {
		return nil, &NotFoundError{id: sha, object: "Commit"}
	}

	repoCommit := &Commit{
		id:      oid,
		Sha:     sha,
		Message: commit.Message(),
		Tree:    &Tree{id: commit.TreeId(), Sha: commit.TreeId().String()},
		Parents: make([]*CommitParent, commit.ParentCount()),
		repo:    r,
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

	return repoCommit, nil
}

// CreateCommit creates a new commit in the repository
func (r *Repo) CreateCommit(treeSha, msg string, parentShas []string) (*Commit, error) {
	treeID, err := git.NewOid(treeSha)
	if err != nil {
		return nil, err
	}

	tree, err := r.repo.LookupTree(treeID)
	if err != nil {
		return nil, &NotFoundError{id: treeSha, object: "Commit Tree"}
	}

	parents := make([]*git.Commit, len(parentShas))
	for i := 0; i < len(parents); i++ {
		oid, err := git.NewOid(parentShas[i])
		if err != nil {
			return nil, err
		}

		commit, err := r.repo.LookupCommit(oid)
		if err != nil {
			return nil, err
		}
		parents[i] = commit
	}

	sig := &git.Signature{
		Name:  r.user.Name(),
		Email: r.user.Email(),
		When:  time.Now(),
	}

	oid, err := r.repo.CreateCommit("", sig, sig, msg, tree, parents...)
	if err != nil {
		return nil, err
	}

	return r.GetCommit(oid.String())
}

// ChangedFiles between two commits
func (c *Commit) ChangedFiles(other *Commit) ([]*FileChange, error) {
	oldTree, err := c.repo.repo.LookupTree(c.Tree.id)
	if err != nil {
		return nil, err
	}
	newTree, err := c.repo.repo.LookupTree(other.Tree.id)
	if err != nil {
		return nil, err
	}

	diff, err := c.repo.repo.DiffTreeToTree(oldTree, newTree, nil)
	if err != nil {
		return nil, err
	}

	changes := []*FileChange{}
	deltas, _ := diff.NumDeltas()
	for i := 0; i < deltas; i++ {
		delta, _ := diff.GetDelta(i)
		if delta.OldFile.Oid == nil || delta.OldFile.Oid.IsZero() {
			changes = append(changes, &FileChange{Path: delta.NewFile.Path, Action: "create"})
		} else if delta.NewFile.Oid == nil || delta.NewFile.Oid.IsZero() {
			changes = append(changes, &FileChange{Path: delta.OldFile.Path, Action: "delete"})
		} else {
			changes = append(changes, &FileChange{Path: delta.NewFile.Path, Action: "update"})
		}
	}

	return changes, nil
}
