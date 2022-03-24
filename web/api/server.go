package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/swarnakumar/go-identity/web/jwt"
	"go.uber.org/zap"
	"net/http"

	"github.com/swarnakumar/go-identity/db"
	"github.com/swarnakumar/go-identity/db/sql/sqlc"
)

type Server interface {
	RenderApiResponse(w http.ResponseWriter, msg any, statusCode int)

	GetRequestUser(r *http.Request) *sqlc.GetUserByEmailRow

	GetDbClient() *db.Client
	GetLogger() *zap.SugaredLogger
	GetJWT() *jwt.JWT
}

func MakeRouter(s Server) chi.Router {
	router := chi.NewRouter()

	router.Post("/generate-token", HandleGenerateToken(s))
	router.Post("/verify-token", HandleVerifyToken(s))

	loginGate := GetLoginGateMiddleware(s)

	router.Group(func(r chi.Router) {
		r.Use(loginGate)
		r.Get("/me", HandleMe(s))
		r.Post("/refresh-token", HandleRefreshToken(s))
		r.Post("/change-password", HandleChangePwd(s))
	})

	return router
}

func MakeJWTRouter(s Server) chi.Router {
	router := chi.NewRouter()

	router.Post("/create", HandleGenerateToken(s))
	router.Post("/verify", HandleVerifyToken(s))

	loginGate := GetLoginGateMiddleware(s)

	router.Group(func(r chi.Router) {
		r.Use(loginGate)
		r.Post("/refresh", HandleRefreshToken(s))
	})

	return router
}
