package web

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"html/template"
	"net/http"

	"github.com/swarnakumar/go-identity/config"
	"github.com/swarnakumar/go-identity/db"
	"github.com/swarnakumar/go-identity/db/sql/sqlc"
	"github.com/swarnakumar/go-identity/web/cookie"
	"github.com/swarnakumar/go-identity/web/jwt"
	"github.com/swarnakumar/go-identity/web/middleware"
	"github.com/swarnakumar/go-identity/web/templates"
)

type Server struct {
	logger *zap.SugaredLogger
	router *chi.Mux
	db     *db.Client
	cookie *cookie.Cookie
	jwt    *jwt.JWT
}

func (s *Server) Close() {
	s.db.Close()
	_ = s.logger.Sync()
}

func New(ctx context.Context) *Server {
	j, err := jwt.New(config.JWTPrivateKeyFile, config.JWTPublicKeyFile, config.JWTAudience, config.JWTIssuer)
	if err != nil {
		panic(err)
	}
	c, err := cookie.New(config.CookieHashKey, config.CookieBlockKey)
	if err != nil {
		panic(err)
	}

	logger, _ := zap.NewProduction()
	sugar := logger.Sugar()

	r := chi.NewRouter()
	client := db.New(ctx)

	server := Server{
		logger: sugar,
		router: r,
		db:     client,
		cookie: c,
		jwt:    j,
	}

	return &server
}

func (s *Server) GetDbClient() *db.Client {
	return s.db
}

func (s *Server) GetLogger() *zap.SugaredLogger {
	return s.logger
}

func (s *Server) GetJWT() *jwt.JWT { return s.jwt }

func (s *Server) GetCookie() *cookie.Cookie { return s.cookie }

func (s *Server) URLParam(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}

func (s *Server) RenderApiResponse(w http.ResponseWriter, msg any, statusCode int) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")

	if msg != nil {

		resp, err := json.Marshal(msg)
		if err != nil {
			s.logger.Errorw("Unable to convert to JSON", "Message", msg, "Error", err)
		}

		_, err = w.Write(resp)
		if err != nil {
			s.logger.Errorw("Unable to write to HttpWriter", "Response", resp, "Error", err)
		}
	}
}

func (s *Server) GetRequestUser(r *http.Request) *sqlc.GetUserByEmailRow {
	return middleware.GetUser(r)
}

type errorParams struct {
	Error string
}

func (s *Server) RenderErrorAlert(w http.ResponseWriter, msg string) {
	m := errorParams{Error: msg}
	err := templates.ErrorAlert.Execute(w, m)
	if err != nil {
		s.logger.Errorw("Unable to execute template",
			"Template", "error-alert", "Message", msg, "Error", err)
	}
}

type successParams struct {
	Message string
}

func (s *Server) RenderSuccessAlert(w http.ResponseWriter, msg string) {
	m := successParams{Message: msg}
	err := templates.SuccessAlert.Execute(w, m)
	if err != nil {
		s.logger.Errorw("Unable to execute template",
			"Template", "success-alert", "Message", msg, "Error", err)
	}
}

func (s *Server) ExecuteTemplate(
	w http.ResponseWriter,
	r *http.Request,
	t *template.Template,
	params *map[string]interface{},
) {
	err := templates.ExecuteTemplate(w, r, t, s.GetRequestUser(r), params)
	if err != nil {
		s.logger.Errorw(
			"Unable to execute template",
			"Template", t.Name(),
			"Params", *params,
			"Error", err)
	}
}
