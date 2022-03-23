package middleware

import (
	"net/http"

	"github.com/swarnakumar/go-identity/web/templates"
)

var adminGateTemplate = templates.Parse("admin-gate.html")

func AdminGateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUser(r)
		if user == nil {
			handleNoUser(w, r)
			return
		}

		if !user.IsAdmin {
			_ = templates.ExecuteTemplate(w, r, adminGateTemplate, user, nil)
			return
		}

		next.ServeHTTP(w, r)
	})
}
