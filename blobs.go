package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/libgit2/git2go"
)

// Blob represents a blob
type Blob struct {
	Sha  string `json:"sha"`
	Size int64  `json:"size"`
}

// BlobCreateParams is the JSON object sent when creating a new blob
type BlobCreateParams struct {
	Content  string `json:"content"`
	Encoding string `json:"encoding"`
}

// CreateBlob uploads a new blob and stores it in the repository object db
func CreateBlob(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var oid *git.Oid
	var reader io.Reader
	if r.Header.Get("Content-Type") == rawContentType {
		reader = r.Body
	} else {
		blobParams := &BlobCreateParams{}
		jsonDecoder := json.NewDecoder(r.Body)
		err := jsonDecoder.Decode(blobParams)
		if err != nil {
			InternalServerError(w, fmt.Sprintf("Could not parse blob params: %v", err))
			return
		}

		if blobParams.Encoding != "base64" {
			InternalServerError(w, fmt.Sprintf("Only base64 encoding supported. Encoding set to: %v", blobParams.Encoding))
			return
		}

		reader = base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(blobParams.Content))
	}

	oid, err := repo.CreateBlobFromChunks("", func(maxLen int) ([]byte, error) {
		b := make([]byte, maxLen)
		_, err := reader.Read(b)
		return b, err
	})

	blob, err := repo.LookupBlob(oid)
	if err != nil {
		InternalServerError(w, fmt.Sprintf("Blob could not be read after creation: %v", err))
		return
	}

	sendJSON(w, 200, &Blob{Sha: oid.String(), Size: blob.Size()})
}

// GetBlob returns information about a blob in the repository.
// If the Content-Type is set to "application/vnd.netlify.raw" it will return
// the actual blob contents
func GetBlob(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	sha := params.ByName("sha")
	oid, err := git.NewOid(sha)
	if err != nil {
		InternalServerError(w, fmt.Sprintf("Invalid sha: %v", err))
		return
	}

	blob, err := repo.LookupBlob(oid)
	if err != nil {
		NotFoundError(w, fmt.Sprintf("Blob not found for %v: %v", sha, err))
		return
	}

	if r.Header.Get("Content-Type") == rawContentType {
		odb, err := repo.Odb()
		if err != nil {
			InternalServerError(w, fmt.Sprintf("Could not open repository backend: %v", err))
			return
		}

		reader, err := odb.NewReadStream(blob.Id())
		if err != nil {
			InternalServerError(w, fmt.Sprintf("Unable to read blob: %v", err))
			return
		}
		defer reader.Close()
		w.Header().Add("Content-Type", "application/octet-stream")
		io.Copy(w, reader)
		return
	}

	sendJSON(w, 200, &Blob{Sha: blob.Id().String(), Size: blob.Size()})
}
