package templates

import (
	"embed"
	"github.com/gorilla/csrf"
	"github.com/swarnakumar/go-identity/db/sql/sqlc"
	"html/template"
	"net/http"
	"net/url"
)

//go:embed *
var files embed.FS

type Crumb struct {
	Text string
	Link string
}

var SuccessAlert = template.Must(template.ParseFS(files, "components/success-alert.html"))
var ErrorAlert = template.Must(template.ParseFS(files, "components/error-alert.html"))

type Number interface {
	int | int8 | int16 | int32 | int64 | float32 | float64
}

var funcMap = template.FuncMap{
	// First Index in Group
	"first": func(index, listLength int) bool {
		return index == 0
	},
	// Last Index in Group
	"last": func(index, listLength int) bool {
		return index == listLength-1
	},
	// Encode String for URL
	"urlEncode": func(s string) string {
		return url.QueryEscape(s)
	},
	// Integer Arithmetic
	"inc": func(v int) int { return v + 1 },
	"dec": func(v int) int { return v - 1 },
}

func Parse(templates ...string) *template.Template {
	templates = append([]string{"layout.html"}, templates...)
	return template.Must(
		template.New("layout.html").Funcs(funcMap).ParseFS(files, templates...))
}

func ParsePartial(name string, templates ...string) *template.Template {
	return template.Must(
		template.New(name).Funcs(funcMap).ParseFS(files, templates...))
}

func ExecuteTemplate(
	w http.ResponseWriter,
	r *http.Request,
	t *template.Template,
	user *sqlc.GetUserByEmailRow,
	params *map[string]interface{},
) error {
	newMap := make(map[string]interface{})
	newMap["currentUser"] = user
	newMap[csrf.TemplateTag] = csrf.TemplateField(r)
	if params != nil {
		for key, value := range *params {
			newMap[key] = value
		}
	}
	err := t.Execute(w, newMap)
	return err
}
