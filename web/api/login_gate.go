package api

import (
	"net/http"

	mw "github.com/swarnakumar/go-identity/web/middleware"
)

func GetLoginGateMiddleware(s Server) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := mw.GetUser(r)
			if user == nil {
				err := ErrorResponse{
					Title:  "invalid-credentials",
					Detail: "Invalid Credentials: Token not passed, or already expired!",
					Status: http.StatusUnauthorized,
				}
				s.RenderApiResponse(w, err, http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
