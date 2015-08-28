package api

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var (
	indexPage = `<!doctype html>
<html>
  <head><title>Local Netlify CMS Backend</title></head>
  <body>
    <h1>Local Netlify CMS Backend</h1>
    <p>
      This is a simple local backend for <a href="https://github.com/netlify/netlify-cms">Netlify's CMS</a>
    </p>
  </body>
</html>
  `
)

// Index serves a basic index page
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, indexPage)
}
