package main

import (
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// GetFile returns information about a file or directory in the repository.
// If the Content-Type is set to "application/vnd.netlify.raw" it will return
// the actual file contents (or an error if a directory)
func GetFile(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	pathname := params.ByName("path")[1:]
	file, err := currentRepo.GetFile(pathname)
	if err != nil {
		HandleError(w, err)
		return
	}

	if file.Type == "dir" {
		sendJSON(w, 200, file.Files)
		return
	}

	if r.Header.Get("Content-Type") != rawContentType {
		sendJSON(w, 200, file)
	}

	blob, err := currentRepo.GetBlob(file.Sha)
	if err != nil {
		HandleError(w, err)
		return
	}
	io.Copy(w, blob)
}
