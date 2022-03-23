package web

import (
	"github.com/go-chi/chi/v5"
	mw "github.com/swarnakumar/go-identity/web/middleware"
	webapp "github.com/swarnakumar/go-identity/web/webapp"
	"github.com/swarnakumar/go-identity/web/webapp/admin"
	"net/http"
	"os"
	"path/filepath"
)

func pong(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

func setupLoginRoutes(s *Server) {
	// All are public routes
	s.router.Get("/login", webapp.RenderLoginPage(s))
	s.router.Post("/login", webapp.HandleLoginCredentials(s))
	s.router.Get("/logout", webapp.HandleLogout(s))
	s.router.Post("/logout", webapp.HandleLogout(s))
}

func getAdminRouter(s *Server) chi.Router {
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

func (s *Server) setupRoutes() {
	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "web", "static")

	s.SetupStatic("/static", filesDir)

	s.router.Get("/ping", pong)

	setupLoginRoutes(s)

	// Behind Login Gate
	s.router.Group(func(r chi.Router) {
		r.Use(mw.LoginGateMiddleware)
		s.router.Get("/", webapp.RenderHomePage(s))
		r.Get("/change-password", webapp.RenderChangePwd(s))
		r.Post("/change-password", webapp.HandleChangePwd(s))
	})

	// Admin
	s.router.Mount("/admin", getAdminRouter(s))

}
