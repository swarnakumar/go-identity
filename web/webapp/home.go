package webapp

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/swarnakumar/go-identity/web/templates"
)

var homeTemplate = templates.Parse("home.html")

func RenderHomePage(s Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := s.GetRequestUser(r)

		if user == nil {
			values := url.Values{}
			values.Add("next", r.URL.String())

			u := fmt.Sprintf("/login?%s", values.Encode())
			http.Redirect(w, r, u, http.StatusFound)
			return
		}

		p1 := map[string]interface{}{"UserDetails": user}
		s.ExecuteTemplate(w, r, homeTemplate, &p1)
	}
}
