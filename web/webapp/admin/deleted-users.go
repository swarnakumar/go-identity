package admin

import (
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/swarnakumar/go-identity/db/sql/sqlc"
	"github.com/swarnakumar/go-identity/web/templates"
	"github.com/swarnakumar/go-identity/web/webapp"
)

var deletedUsersTemplate = templates.Parse("admin/user-deletions.html", "components/pagination.html")

func RenderDeletedUsers(s webapp.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user := s.GetRequestUser(r)
		if user == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		page, err := strconv.ParseInt(r.URL.Query().Get("page"), 10, 32)
		if err != nil {
			page = 1
		}

		limit, err := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 32)
		if err != nil {
			limit = 100
		}

		db := s.GetDbClient()

		c1 := make(chan []*sqlc.UserDeletion)
		c2 := make(chan int64)

		go func() {

			offset := int32((page - 1) * limit)
			listing, err := db.Users.GetDeletions(ctx, offset, int32(limit))
			if err != nil {
				s.GetLogger().Errorw("Unable to get user list", "Error", err)
			}
			c1 <- listing
		}()

		go func() {
			count := db.Users.GetDeletionsCount(ctx)
			c2 <- count
		}()

		var listing []*sqlc.UserDeletion
		var count int64

		for i := 0; i < 2; i++ {
			select {
			case res := <-c1:
				listing = res
			case c := <-c2:
				count = c
			}
		}

		pageCount := int64(math.Ceil(float64(count) / float64(limit)))

		if page > pageCount {
			http.Redirect(w, r, "/admin/users", http.StatusFound)
			return
		}

		crumbs := []templates.Crumb{
			{Text: "Home", Link: "/"},
			{Text: "Admin", Link: "/admin"},
			{Text: "Users", Link: "/admin/users"},
		}

		if page >= 2 {
			c := templates.Crumb{
				Text: fmt.Sprintf("Page %d", page),
				Link: fmt.Sprintf("/admin/users?page=%d", page),
			}
			crumbs = append(crumbs, c)
		}

		params := map[string]interface{}{
			"Users":     listing,
			"CurrPage":  page,
			"PageCount": pageCount,
			"crumbs":    crumbs,
			"Link":      "/admin/users",
		}

		s.ExecuteTemplate(w, r, deletedUsersTemplate, &params)
	}
}
