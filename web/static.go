package web

import (
	"embed"
	"io/fs"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
)

//go:embed static
var embeddedFiles embed.FS

func (s *Server) SetupStatic(staticPath, staticRoot string) {
	if strings.ContainsAny(staticPath, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if staticPath != "/" && staticPath[len(staticPath)-1] != '/' {
		s.router.Get(staticPath, http.RedirectHandler(staticPath+"/", http.StatusMovedPermanently).ServeHTTP)
		staticPath += "/"
	}
	staticPath += "*"

	s.router.Get(staticPath, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")

		env := os.Getenv("ENV")

		var fsys fs.FS
		if env == "dev" || env == "test" {
			fsys = os.DirFS(staticRoot)
		} else {
			fsys, _ = fs.Sub(embeddedFiles, staticRoot)
		}

		fileServer := http.StripPrefix(pathPrefix, http.FileServer(http.FS(fsys)))
		fileServer.ServeHTTP(w, r)
	})
}
