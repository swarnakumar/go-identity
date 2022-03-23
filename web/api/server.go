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

func MakeApiRouter(s Server) chi.Router {
	apiRouter := chi.NewRouter()

	apiRouter.Post("/generate-token", HandleGenerateToken(s))
	apiRouter.Post("/refresh-token", HandleRefreshToken(s))
	apiRouter.Post("/verify-token", HandleVerifyToken(s))
	apiRouter.Post("/change-password", HandleChangePwd(s))

	return apiRouter

}
