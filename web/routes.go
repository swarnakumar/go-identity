package web

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/swarnakumar/go-identity/web/api"
	"github.com/swarnakumar/go-identity/web/server"
	webRouter "github.com/swarnakumar/go-identity/web/webapp/router"
)

func pong(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

func setupRoutes(s *server.Server) {
	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "web", "static")

	r := s.GetRouter()
	setupStatic(r, "/static", filesDir)

	r.Get("/ping", pong)

	w := webRouter.MakeRouter(s)
	r.Mount("/", w)

	apiRouter := api.MakeRouter(s)
	r.Mount("/api", apiRouter)

	jwtRouter := api.MakeJWTRouter(s)
	r.Mount("/jwt", jwtRouter)
}
