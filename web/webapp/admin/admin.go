package admin

import (
	"github.com/swarnakumar/go-identity/web/templates"
	"github.com/swarnakumar/go-identity/web/webapp"
	"net/http"
)

var adminLandingTemplate = templates.Parse("admin/admin-landing.html")

func RenderAdminPage(s webapp.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		crumbs := []templates.Crumb{
			{Text: "Home", Link: "/"},
			{Text: "Admin", Link: "/admin"},
		}
		params := map[string]interface{}{"crumbs": crumbs}
		s.ExecuteTemplate(w, r, adminLandingTemplate, &params)
	}
}
