package repo

import (
	"fmt"
	"strconv"

	"gopkg.in/libgit2/git2go.v22"
)

// Tree a tree in the repository object db
type Tree struct {
	id   *git.Oid
	Sha  string       `json:"sha"`
	Tree []*TreeEntry `json:"tree,omitempty"`
}

// TreeEntry a single entry in a RepoTree
type TreeEntry struct {
	id   *git.Oid
	Path string `json:"path"`
	Type string `json:"type"`
	Mode string `json:"mode"`
	Size int64  `json:"size"`
	Sha  string `json:"sha"`
}

func (r *Repo) newTreeEntry(entry *git.TreeEntry) *TreeEntry {
	var size int64
	var objType string
	if entry.Type == git.ObjectBlob {
		blob, _ := r.repo.LookupBlob(entry.Id)
		size = blob.Size()
		objType = "blob"
	} else {
		objType = "tree"
	}
	return &TreeEntry{
		id:   entry.Id,
		Path: entry.Name,
		Size: size,
		Mode: fmt.Sprintf("%v", entry.Filemode),
		Sha:  entry.Id.String(),
		Type: objType,
	}
}

// GetTree returns a tree from the repository from a sha
func (r *Repo) GetTree(sha string) (*Tree, error) {
	oid, err := git.NewOid(sha)
	if err != nil {
		return nil, err
	}

	tree, err := r.repo.LookupTree(oid)
	if err != nil {
		return nil, &NotFoundError{id: sha, object: "Tree"}
	}

	repoTree := &Tree{
		id:   oid,
		Sha:  sha,
		Tree: make([]*TreeEntry, tree.EntryCount()),
	}

	var i uint64
	for i = 0; i < tree.EntryCount(); i++ {
		repoTree.Tree[i] = r.newTreeEntry(tree.EntryByIndex(i))
	}

	return repoTree, nil
}

// CreateTree creates a new tree in the repo.
// If baseSha is not empty, it will be based on an existing tree
func (r *Repo) CreateTree(baseSha string, entries []*TreeEntry) (*Tree, error) {
	builder, err := r.repo.TreeBuilder()
	if err != nil {
		return nil, err
	}
	defer builder.Free()

	if baseSha != "" {
		baseID, err := git.NewOid(baseSha)
		if err != nil {
			return nil, err
		}

		base, err := r.repo.LookupTree(baseID)
		if err != nil {
			return nil, &NotFoundError{id: baseSha, object: "Base Tree"}
		}

		var i uint64
		for i = 0; i < base.EntryCount(); i++ {
			entry := base.EntryByIndex(i)
			err = builder.Insert(entry.Name, entry.Id, int(entry.Filemode))
			if err != nil {
				return nil, err
			}
		}
	}

	for _, entry := range entries {
		oid, err := git.NewOid(entry.Sha)
		if err != nil {
			return nil, err
		}
		mode, err := strconv.Atoi(entry.Mode)
		if err != nil {
			return nil, err
		}
		err = builder.Insert(entry.Path, oid, mode)
		if err != nil {
			return nil, err
		}
	}

	oid, err := builder.Write()
	if err != nil {
		return nil, err
	}

	return r.GetTree(oid.String())
}
