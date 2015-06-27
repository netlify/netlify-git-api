package repo

import (
	"io"

	"gopkg.in/libgit2/git2go.v22"
)

// Blob represents a blob
type Blob struct {
	Sha  string `json:"sha"`
	Size int64  `json:"size"`
	id   *git.Oid
	repo *Repo
	i    int64  // current reading index
	data []byte // data
}

func (b *Blob) Read(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}
	if b.i == 0 {
		blob, err := b.repo.repo.LookupBlob(b.id)
		if err != nil {
			return 0, err
		}
		b.data = blob.Contents()
	}
	if b.i >= int64(len(b.data)) {
		return 0, io.EOF
	}
	n = copy(p, b.data[b.i:])
	b.i += int64(n)
	return
}

// WriteTo writes the content of a blob to a writer
// Seems NewReadStream is not supported in the libgit2 backend :/
// func (b *Blob) WriteTo(w io.Writer) (n int64, err error) {
// 	odb, err := b.repo.repo.Odb()
// 	if err != nil {
// 		return n, err
// 	}
//
// 	log.Println("Opening readstream from odb")
// 	reader, err := odb.NewReadStream(b.id)
// 	if err != nil {
// 		log.Printf("Error opening stream: %v", err)
// 		return n, err
// 	}
// 	log.Printf("Copying %v to %v for %v\n", reader, w, b.Sha)
// 	return io.Copy(w, reader)
// }

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
		id:   oid,
		repo: r,
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
