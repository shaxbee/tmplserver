package tmplserver // import "github.com/shaxbee/tmplserver"

import (
	"net/http"
	"os"
	"path"
	"strings"

	"log"
)

type tmplServer struct {
	resolvers []Resolver
}

// NewServer creates server that uses resolvers to locate resources
func NewServer(resolvers ...Resolver) http.Handler {
	return &tmplServer{resolvers}
}

func (s *tmplServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := r.Body.Close(); err != nil {
		http.Error(w, "408 Request timeout", http.StatusRequestTimeout)
		return
	}

	if strings.Contains(r.URL.Path, "..") {
		http.Error(w, "400 Invalid URL Path", http.StatusBadRequest)
		return
	}

	if r.URL.Path == "/" {
		r.URL.Path = "/index.html"
	}

	r.URL.Path = path.Clean(r.URL.Path)

	for _, resolver := range s.resolvers {
		mod, tr, err := resolver.Resolve(r.URL.Path[1:])
		switch {
		case os.IsPermission(err):
			log.Print(err)
			http.Error(w, "403 Forbidden", http.StatusForbidden)
			return
		case err != nil:
			log.Print(err)
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
			return
		case tr != nil:
			http.ServeContent(w, r, r.URL.Path, mod, tr)
			return
		}
	}

	http.Error(w, "404 Resource Not Found", http.StatusNotFound)
}
