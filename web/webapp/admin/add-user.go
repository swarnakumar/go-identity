package admin

import (
	"net/http"
	"strconv"

	"github.com/swarnakumar/go-identity/web/webapp"
)

func HandleAddNewUser(s webapp.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			s.GetLogger().Errorw("Unable to parse login credentials form!", "Error", err)
			return
		}
		ctx := r.Context()

		currentUser := s.GetRequestUser(r)

		email := r.FormValue("email")
		password := r.FormValue("password")
		isAdmin, _ := strconv.ParseBool(r.FormValue("is-admin"))

		db := s.GetDbClient()
		_, err = db.Users.Create(ctx, email, password, isAdmin, &(currentUser.Email))

		if err == nil {
			s.RenderSuccessAlert(w, "Successfully Created User!!!")
		} else {
			s.RenderErrorAlert(w, err.Error())
		}

	}
}
