package server

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/gorilla/csrf"

	"github.com/swarnakumar/go-identity/config"
	mw "github.com/swarnakumar/go-identity/web/middleware"
)

func (s *Server) InitRouter() {

	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	corsMiddleware := cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	csrfMiddleware := csrf.Protect(
		[]byte(config.CSRFAuthKey),
		csrf.Secure(config.UseHttps),
	)

	jwtVerifier := s.jwt.GetVerifierMiddleware()

	authMiddleware := mw.GetAuthMiddleware(s.db)

	s.router.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,

		csrfMiddleware,
		corsMiddleware,
		jwtVerifier,

		authMiddleware,
	)
}
