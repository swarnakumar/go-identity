package api

import "net/http"

func HandleMe(s Server) http.HandlerFunc {
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

		s.RenderApiResponse(w, currUser, http.StatusOK)
	}
}
