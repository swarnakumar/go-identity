package admin

import (
	"net/http"
	"net/url"

	"github.com/swarnakumar/go-identity/web/webapp"
)

func HandleUserDeletion(s webapp.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email, _ := url.QueryUnescape(s.URLParam(r, "email"))

		db := s.GetDbClient()
		currentUser := s.GetRequestUser(r)
		err := db.Users.Delete(r.Context(), email, &currentUser.Email)
		if err == nil {
			w.Header().Set("HX-Redirect", "/admin/users")
		} else {
			s.RenderErrorAlert(w, err.Error())
		}
	}
}
