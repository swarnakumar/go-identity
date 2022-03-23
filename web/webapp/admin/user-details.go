package admin

import (
	"net/http"
	"net/url"

	"github.com/swarnakumar/go-identity/web/templates"
	"github.com/swarnakumar/go-identity/web/webapp"
)

var userDetailsTemplate = templates.Parse("admin/user-details.html")

func RenderUserDetails(s webapp.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email, _ := url.QueryUnescape(s.URLParam(r, "email"))

		db := s.GetDbClient()
		user, err := db.Users.GetByEmail(r.Context(), email)

		crumbs := []templates.Crumb{
			{Text: "Home", Link: "/"},
			{Text: "Admin", Link: "/admin"},
			{Text: "Users", Link: "/admin/users"},
			{Text: email, Link: r.URL.String()},
		}

		params := map[string]interface{}{"crumbs": crumbs}

		if err == nil {
			params["User"] = user
		}
		s.ExecuteTemplate(w, r, userDetailsTemplate, &params)

	}
}

var userChangeResultTemplate = templates.ParsePartial("change-user.html", "admin/change-user.html")

func HandleUserChange(s webapp.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email, err := url.QueryUnescape(s.URLParam(r, "email"))
		if email == "" || err != nil {
			s.RenderErrorAlert(w, "Unable to find User")
		}

		isAdmin := r.FormValue("is-admin") == "on"
		isActive := r.FormValue("is-active") == "on"

		currUser := s.GetRequestUser(r).Email
		db := s.GetDbClient()
		user, err := db.Users.UpdateAttributes(r.Context(), email, isAdmin, isActive, &currUser)

		if err == nil {
			params := map[string]interface{}{"User": user, "Message": "Update Successful"}
			s.RenderSuccessAlert(w, "Successfully Updated")
			s.ExecuteTemplate(w, r, userChangeResultTemplate, &params)
		} else {
			s.RenderErrorAlert(w, "Nothing to Update!!!")
		}

	}
}

func HandlePwdChange(s webapp.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email, err := url.QueryUnescape(s.URLParam(r, "email"))
		if email == "" || err != nil {
			s.RenderErrorAlert(w, "Unable to find User")
		}

		pwd1 := r.FormValue("pwd1")
		pwd2 := r.FormValue("pwd2")

		if pwd1 != pwd2 || pwd1 == "" {
			s.RenderErrorAlert(w, "The two passwords dont match!!!")
			return
		}

		currUser := s.GetRequestUser(r).Email
		db := s.GetDbClient()
		user, err := db.Users.ChangePassword(r.Context(), email, pwd1, &currUser)

		if err == nil {
			params := map[string]interface{}{"User": user, "Message": "Update Successful"}
			s.RenderSuccessAlert(w, "Successfully Updated")
			s.ExecuteTemplate(w, r, userChangeResultTemplate, &params)
		} else {
			s.RenderErrorAlert(w, "Your password's not complex enough!")
		}

	}
}
