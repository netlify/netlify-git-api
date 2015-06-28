package repo

import (
	"fmt"
	"path"

	"gopkg.in/libgit2/git2go.v22"
)

// File a file in the repository
type File struct {
	id    *git.Oid
	Name  string  `json:"name"`
	Path  string  `json:"path"`
	Size  int64   `json:"size"`
	Sha   string  `json:"sha"`
	Type  string  `json:"type"`
	Files []*File `json:"files,omitempty"`
}

func (r *Repo) newRepoFile(entry *git.TreeEntry, dir string, expand bool) (*File, error) {
	file := &File{
		id:   entry.Id,
		Name: entry.Name,
		Path: path.Join(dir, entry.Name),
		Sha:  entry.Id.String(),
	}

	if entry.Type == git.ObjectBlob {
		blob, err := r.repo.LookupBlob(entry.Id)
		if err != nil {
			return nil, err
		}
		file.Size = blob.Size()
		file.Type = "file"
	} else {
		file.Type = "dir"
	}

	if expand && entry.Type == git.ObjectTree {
		tree, err := r.repo.LookupTree(entry.Id)
		if err != nil {
			return nil, err
		}

		file.Files = make([]*File, tree.EntryCount())

		var i uint64
		for i = 0; i < tree.EntryCount(); i++ {
			entry := tree.EntryByIndex(i)
			newFile, _ := r.newRepoFile(entry, file.Path, false)
			file.Files[i] = newFile
		}
	}

	return file, nil
}

// GetFile finds a file or directory
func (r *Repo) GetFile(pathname string) (*File, error) {
	var entry *git.TreeEntry
	ref, err := r.repo.Head()
	if err != nil {
		return nil, err
	}

	obj, err := r.repo.Lookup(ref.Target())
	if err != nil {
		return nil, err
	}

	commit, ok := obj.(*git.Commit)
	if !ok {
		return nil, fmt.Errorf("Head is not a commit: %v", obj.Type())
	}

	tree, err := commit.Tree()
	if err != nil {
		return nil, err
	}

	if pathname == "" {
		entry = &git.TreeEntry{Name: "", Id: tree.Id(), Type: git.ObjectTree}
	} else {
		entry, err = tree.EntryByPath(pathname)
	}

	if err != nil {
		return nil, &NotFoundError{id: pathname, object: "File or Dir"}
	}

	file, err := r.newRepoFile(entry, path.Dir(pathname), true)
	if err != nil {
		return nil, err
	}
	return file, nil
}
