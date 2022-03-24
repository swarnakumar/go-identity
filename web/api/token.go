package api

import (
	"encoding/json"
	"net/http"
)

type getTokenParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// HandleGenerateToken returns JWT for user
func HandleGenerateToken(s Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := getTokenParams{}
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil || data.Email == "" || data.Password == "" {
			e := ErrorResponse{
				Title:  "invalid-request",
				Detail: "Unable to Parse to JSON",
				Status: http.StatusBadRequest,
			}
			s.RenderApiResponse(w, e, http.StatusBadRequest)
			return
		}

		db := s.GetDbClient()

		ok := db.Users.VerifyPassword(r.Context(), data.Email, data.Password)
		if !ok {
			err := ErrorResponse{
				Title:  "invalid-credentials",
				Detail: "Invalid Credentials: email and/or password are wrong!",
				Status: http.StatusUnauthorized,
			}
			s.RenderApiResponse(w, err, http.StatusUnauthorized)
		} else {
			jwt := s.GetJWT()
			jwtClaims := map[string]interface{}{"user": data.Email}
			token := jwt.Create(jwtClaims)

			resp := map[string]string{"token": token}
			s.RenderApiResponse(w, resp, http.StatusCreated)
		}
	}
}

// HandleRefreshToken generates a new token for user with new expiry.
func HandleRefreshToken(s Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currUser := s.GetRequestUser(r)
		if currUser == nil {
			err := ErrorResponse{
				Title:  "invalid-credentials",
				Detail: "Invalid Credentials: Token not passed, or already expired!",
				Status: http.StatusUnauthorized,
			}
			s.RenderApiResponse(w, err, http.StatusUnauthorized)
			return
		}

		jwtClaims := map[string]interface{}{"user": currUser.Email}
		token := s.GetJWT().Create(jwtClaims)

		resp := map[string]string{"token": token}
		s.RenderApiResponse(w, resp, http.StatusCreated)
	}
}

type verifyTokenParams struct {
	Token string `json:"token"`
}

// HandleVerifyToken verifies token.
func HandleVerifyToken(s Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleInvalidRequest := func() {
			e := ErrorResponse{
				Title:  "invalid-data",
				Detail: "Invalid Data Sent.",
				Status: http.StatusBadRequest,
			}
			s.RenderApiResponse(w, e, http.StatusBadRequest)
		}

		if r.Body == nil {
			handleInvalidRequest()
			return
		}

		decoder := json.NewDecoder(r.Body)
		var t verifyTokenParams
		err := decoder.Decode(&t)
		if err != nil || t.Token == "" {
			handleInvalidRequest()
			return
		}

		handleValidRequest := func(valid bool) {
			resp := map[string]bool{"valid": valid}
			s.RenderApiResponse(w, resp, http.StatusOK)
		}

		j := s.GetJWT()
		token, err := j.Get(t.Token)
		if err != nil || token == nil {
			handleValidRequest(false)
			return
		}

		user, present := token.PrivateClaims()["user"]
		if !present || user == nil || user == "" {
			handleValidRequest(false)
			return
		}

		_, err = s.GetDbClient().Users.GetByEmail(r.Context(), user.(string))
		handleValidRequest(err == nil)
	}
}

//func HandleDeleteToken(w http.ResponseWriter, r *http.Request, s *server.Server) {
//	s.DeleteCookie(w, "jwt")
//	resp := map[string]string{}
//	s.RenderApiResponse(w, resp, http.StatusCreated)
//}
