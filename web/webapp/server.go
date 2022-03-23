package webapp

import (
	"github.com/swarnakumar/go-identity/web/cookie"
	"github.com/swarnakumar/go-identity/web/jwt"
	"go.uber.org/zap"
	"html/template"
	"net/http"

	"github.com/swarnakumar/go-identity/db"
	"github.com/swarnakumar/go-identity/db/sql/sqlc"
)

type Server interface {
	GetRequestUser(r *http.Request) *sqlc.GetUserByEmailRow

	RenderErrorAlert(w http.ResponseWriter, msg string)
	RenderSuccessAlert(w http.ResponseWriter, msg string)

	ExecuteTemplate(
		w http.ResponseWriter,
		r *http.Request,
		t *template.Template,
		params *map[string]interface{},
	)

	GetDbClient() *db.Client
	GetLogger() *zap.SugaredLogger
	GetJWT() *jwt.JWT
	GetCookie() *cookie.Cookie

	URLParam(r *http.Request, key string) string
}
