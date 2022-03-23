package middleware

import (
	"fmt"
	"net/http"
	"net/url"
)

func handleNoUser(w http.ResponseWriter, r *http.Request) {
	values := url.Values{}
	values.Add("next", r.URL.String())

	u := fmt.Sprintf("/login?%s", values.Encode())
	http.Redirect(w, r, u, http.StatusFound)

}

func LoginGateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUser(r)
		if user == nil {
			handleNoUser(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}
