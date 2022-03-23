package web

import (
	"github.com/swarnakumar/go-identity/web/api"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"

	mw "github.com/swarnakumar/go-identity/web/middleware"
	"github.com/swarnakumar/go-identity/web/server"
	"github.com/swarnakumar/go-identity/web/webapp"
	"github.com/swarnakumar/go-identity/web/webapp/admin"
)

func pong(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

func setupLoginRoutes(s *server.Server) {
	// All are public routes
	r := s.GetRouter()
	r.Get("/login", webapp.RenderLoginPage(s))
	r.Post("/login", webapp.HandleLoginCredentials(s))
	r.Get("/logout", webapp.HandleLogout(s))
	r.Post("/logout", webapp.HandleLogout(s))
}

func getAdminRouter(s *server.Server) chi.Router {
	adminRouter := chi.NewRouter()
	adminRouter.Use(mw.AdminGateMiddleware)

	adminRouter.Get("/", admin.RenderAdminPage(s))
	adminRouter.Get("/users", admin.RenderUserListing(s))
	adminRouter.Get("/user-deletions", admin.RenderDeletedUsers(s))
	adminRouter.Post("/users/add", admin.HandleAddNewUser(s))

	adminRouter.Route("/users/{email}", func(r chi.Router) {
		r.Get("/", admin.RenderUserDetails(s))
		r.Get("/changes", admin.RenderChanges(s))
		r.Post("/change", admin.HandleUserChange(s))
		r.Post("/change-password", admin.HandlePwdChange(s))
		r.Post("/delete", admin.HandleUserDeletion(s))
	})

	return adminRouter
}

func setupRoutes(s *server.Server) {
	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "web", "static")

	gr := s.GetRouter()
	setupStatic(gr, "/static", filesDir)

	gr.Get("/ping", pong)

	setupLoginRoutes(s)

	// Behind Login Gate
	gr.Group(func(r chi.Router) {
		r.Use(mw.LoginGateMiddleware)
		r.Get("/", webapp.RenderHomePage(s))
		r.Get("/change-password", webapp.RenderChangePwd(s))
		r.Post("/change-password", webapp.HandleChangePwd(s))
	})

	// Admin
	gr.Mount("/admin", getAdminRouter(s))

	apiRouter := api.MakeApiRouter(s)
	gr.Mount("/api", apiRouter)
}
