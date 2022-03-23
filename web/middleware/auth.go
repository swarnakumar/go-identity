package middleware

import (
	"context"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/swarnakumar/go-identity/db"
	"github.com/swarnakumar/go-identity/db/sql/sqlc"
)

type authUser string

const authUserKey authUser = "user"

func GetAuthMiddleware(db *db.Client) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t, _, err := jwtauth.FromContext(r.Context())
			if t != nil && err == nil {
				user, _ := t.PrivateClaims()["user"]
				if user == "" || user == nil {
					next.ServeHTTP(w, r)
				} else {
					u, err := db.Users.GetByEmail(r.Context(), user.(string))
					if err == nil {
						ctx := context.WithValue(r.Context(), authUserKey, *u)
						next.ServeHTTP(w, r.WithContext(ctx))
					} else {
						next.ServeHTTP(w, r)
					}

				}
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}
}

func GetUser(r *http.Request) *sqlc.GetUserByEmailRow {
	ctx := r.Context()
	user := ctx.Value(authUserKey)
	if user == nil {
		return nil
	}

	u := user.(sqlc.GetUserByEmailRow)
	return &u
}
