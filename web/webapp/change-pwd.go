package webapp

import (
	"net/http"

	"github.com/swarnakumar/go-identity/web/templates"
)

var changePwdTemplate = templates.Parse("change-pwd.html")

func RenderChangePwd(s Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		crumbs := []templates.Crumb{
			{Text: "Home", Link: "/"},
			{Text: "Change Password", Link: r.URL.String()},
		}

		params := map[string]interface{}{"crumbs": crumbs}
		s.ExecuteTemplate(w, r, changePwdTemplate, &params)
	}
}

func HandleChangePwd(s Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		oldPwd := r.FormValue("current")
		currUser := s.GetRequestUser(r).Email

		db := s.GetDbClient()
		oldIsCorrect := db.Users.VerifyPassword(r.Context(), currUser, oldPwd)
		if !oldIsCorrect {
			s.RenderErrorAlert(w, "Your current password is wrong!!!")
			return
		}

		pwd1 := r.FormValue("pwd1")
		pwd2 := r.FormValue("pwd2")

		if pwd1 != pwd2 || pwd1 == "" {
			s.RenderErrorAlert(w, "The two passwords dont match!!!")
			return
		}

		_, err := db.Users.ChangePassword(r.Context(), currUser, pwd1, &currUser)

		if err == nil {
			s.RenderSuccessAlert(w, "Successfully Updated!")
		} else {
			s.RenderErrorAlert(w, "Your password's not complex enough!")
		}

	}
}
