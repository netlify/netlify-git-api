package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
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
func CreateBlob(w http.ResponseWriter, r *http.Request, params httprouter.Params, ctx context.Context) {
	currentRepo := getRepo(ctx)

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
	blob, err := currentRepo.PutBlob(reader)
	if err != nil {
		HandleError(w, err)
		return
	}

	sendJSON(w, 200, blob)
}

// GetBlob returns information about a blob in the repository.
// If the Content-Type is set to "application/vnd.netlify.raw" it will return
// the actual blob contents
func GetBlob(w http.ResponseWriter, r *http.Request, params httprouter.Params, ctx context.Context) {
	currentRepo := getRepo(ctx)
	sha := params.ByName("sha")
	blob, err := currentRepo.GetBlob(sha)
	if err != nil {
		HandleError(w, err)
		return
	}

	if r.Header.Get("Content-Type") == rawContentType {
		w.Header().Add("Content-Type", "application/octet-stream")
		io.Copy(w, blob)
		return
	}

	sendJSON(w, 200, blob)
}
