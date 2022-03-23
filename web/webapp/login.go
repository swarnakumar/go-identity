package webapp

import (
	"fmt"
	"net/http"

	"github.com/swarnakumar/go-identity/web/templates"
)

var loginTemplate = templates.Parse("login/login.html")

func RenderLoginPage(s Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := s.GetRequestUser(r)
		fmt.Println("LINE 14:", user)
		if user != nil {
			next := r.URL.Query().Get("next")
			if next == "" {
				next = "/"
			}
			fmt.Println("LINE 20:", next)
			http.Redirect(w, r, next, http.StatusFound)
			return
		}

		s.ExecuteTemplate(w, r, loginTemplate, nil)
	}
}

func HandleLoginCredentials(s Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			s.GetLogger().Errorw("Unable to parse login credentials form!", "Error", err)
			return
		}

		db := s.GetDbClient()

		email := r.FormValue("email")
		password := r.FormValue("password")
		valid := db.Users.VerifyPassword(r.Context(), email, password)

		if !valid {
			s.RenderErrorAlert(w, "Invalid Credentials")
			return
		}

		jwtClaims := map[string]interface{}{"user": email}
		s.GetJWT().Set(w, jwtClaims)

		next := r.URL.Query().Get("next")
		if next == "" {
			next = "/"
		}
		w.Header().Set("HX-Redirect", "/")

	}
}

func HandleLogout(s Server) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie := s.GetCookie()
		cookie.Delete(w, "jwt")
		http.Redirect(w, r, "/", http.StatusFound)
	}
}
