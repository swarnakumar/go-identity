package middleware

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/swarnakumar/go-identity/config"
	jwt2 "github.com/swarnakumar/go-identity/web/jwt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/swarnakumar/go-identity/db"
)

var email = "authmiddleware-testing@xyz.com"
var pwd = "qwerty1234@96219authmiddleware"

func TestMain(m *testing.M) {
	ctx := context.Background()

	dbClient := db.New(ctx)
	defer dbClient.Close()

	_, _ = dbClient.Users.Create(ctx, email, pwd, false, nil)
	m.Run()
}

func TestGetAuthMiddleware(t *testing.T) {
	ctx := context.Background()

	dbClient := db.New(ctx)
	defer dbClient.Close()

	jwt, _ := jwt2.New(config.JWTPrivateKeyFile, config.JWTPublicKeyFile, config.JWTAudience, config.JWTIssuer)

	claims := map[string]interface{}{"user": email}
	token := jwt.Create(claims)

	r := chi.NewRouter()

	nextHandler := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := GetUser(r)
			assert.NotNil(t, user)
			assert.Equal(t, email, user.Email)
		})
	}

	r.Use(jwt.GetVerifierMiddleware(), GetAuthMiddleware(dbClient))
	r.Use(nextHandler)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {

	})

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	r.ServeHTTP(rr, req)

	// User that doesnt exist!!!
	claims = map[string]interface{}{"user": "ance jkdcbnkejbc"}
	token = jwt.Create(claims)

	r = chi.NewRouter()

	nextHandler = func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := GetUser(r)
			assert.Nil(t, user)
		})
	}

	r.Use(jwt.GetVerifierMiddleware(), GetAuthMiddleware(dbClient))
	r.Use(nextHandler)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {

	})

	rr = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/", nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	r.ServeHTTP(rr, req)

	// Invalid JWT. User should again be null
	token = " ejkncejkveiwhviv"
	rr = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/", nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	r.ServeHTTP(rr, req)

}
