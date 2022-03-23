package api

import (
	"encoding/json"
	"net/http"
)

type changePwdParams struct {
	CurrentPwd string `json:"current_pwd"`
	NewPwd     string `json:"new_pwd"`
}

// HandleChangePwd changes password for currently logged-in user
func HandleChangePwd(s Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// swagger:route POST /change-password password
		//
		// Changes password for currently logged in user.
		//
		//	 	Consumes:
		//		- application/json
		//
		//		Produces:
		//		- application/json
		//
		//		Schemes: http, https
		//
		//		Security:
		//			api_key:
		//
		//		Parameters:
		//		+ name: current_pwd
		//		  in: body
		//		  description: The user's current password
		//		  required: true
		//		  format: string
		//		+ name: new_pwd
		//		  in: body
		//		  description: User's new password
		//		  required: true
		//		  format: string
		//
		//		Responses:
		//			200: Successfully Changed
		//			401: Invalid Credentials
		//			400: Invalid Request
		data := changePwdParams{}
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			e := ErrorResponse{
				Title:  "invalid-request",
				Detail: "Unable to Parse to JSON",
				Status: http.StatusBadRequest,
			}
			s.RenderApiResponse(w, e, http.StatusBadRequest)
			return
		}

		currUser := s.GetRequestUser(r).Email

		ctx := r.Context()
		db := s.GetDbClient()

		oldIsCorrect := db.Users.VerifyPassword(ctx, currUser, data.CurrentPwd)
		if !oldIsCorrect {
			err := ErrorResponse{
				Title:  "invalid-credentials",
				Detail: "Invalid Credentials to Change Password!",
				Status: http.StatusUnauthorized,
			}

			s.RenderApiResponse(w, err, http.StatusUnauthorized)
			return
		}

		_, err = db.Users.ChangePassword(ctx, currUser, data.NewPwd, &currUser)
		if err == nil {
			s.RenderApiResponse(w, nil, http.StatusCreated)
		} else {
			e := ErrorResponse{
				Title:  "invalid-credentials",
				Detail: err.Error(),
				Status: http.StatusBadRequest,
			}
			s.RenderApiResponse(w, e, http.StatusBadRequest)
		}
	}
}
