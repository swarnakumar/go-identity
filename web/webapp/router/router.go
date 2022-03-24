package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/swarnakumar/go-identity/config"
	mw "github.com/swarnakumar/go-identity/web/middleware"
	"github.com/swarnakumar/go-identity/web/webapp"
	"github.com/swarnakumar/go-identity/web/webapp/admin"
)

func getAdminRouter(s webapp.Server) chi.Router {
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

func MakeRouter(s webapp.Server) chi.Router {
	router := chi.NewRouter()

	var csrfMiddleware = csrf.Protect(
		[]byte(config.CSRFAuthKey),
		csrf.Secure(config.UseHttps),
	)

	router.Use(csrfMiddleware)

	// Login routes
	router.Get("/login", webapp.RenderLoginPage(s))
	router.Post("/login", webapp.HandleLoginCredentials(s))
	router.Get("/logout", webapp.HandleLogout(s))
	router.Post("/logout", webapp.HandleLogout(s))

	// Behind Login Gate
	router.Group(func(r chi.Router) {
		r.Use(csrfMiddleware, mw.LoginGateMiddleware)
		r.Get("/", webapp.RenderHomePage(s))
		r.Get("/change-password", webapp.RenderChangePwd(s))
		r.Post("/change-password", webapp.HandleChangePwd(s))
	})

	router.Mount("/admin", getAdminRouter(s))

	return router
}
