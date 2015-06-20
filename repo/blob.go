package repo

import (
	"bytes"
	"io"

	"github.com/libgit2/git2go"
)

// Blob represents a blob
type Blob struct {
	Sha  string `json:"sha"`
	Size int64  `json:"size"`
	id   *git.Oid
	repo Repo
}

func (b *Blob) Read(p []byte) (n int, err error) {
	buf := bytes.NewBuffer(p)
	i, err := b.WriteTo(buf)
	return int(i), err
}

// WriteTo writes the content of a blob to a writer
func (b *Blob) WriteTo(w io.Writer) (n int64, err error) {
	odb, err := b.repo.repo.Odb()
	if err != nil {
		return n, err
	}

	reader, err := odb.NewReadStream(b.id)
	if err != nil {
		return n, err
	}
	defer reader.Free()
	return io.Copy(w, reader)
}

// GetBlob returns a single blob from the repo
func (r *Repo) GetBlob(sha string) (*Blob, error) {
	oid, err := git.NewOid(sha)
	if err != nil {
		return nil, err
	}

	blob, err := r.repo.LookupBlob(oid)
	if err != nil {
		return nil, &NotFoundError{sha, "Blob"}
	}

	return &Blob{
		Sha:  sha,
		Size: blob.Size(),
	}, nil
}

// PutBlob writes a new blob to the repo and returns the new sha
func (r *Repo) PutBlob(reader io.Reader) (*Blob, error) {
	oid, err := r.repo.CreateBlobFromChunks("", func(maxLen int) ([]byte, error) {
		b := make([]byte, maxLen)
		l, err := reader.Read(b)
		return b[0:l], err
	})
	if err != nil {
		return nil, err
	}

	return r.GetBlob(oid.String())
}
