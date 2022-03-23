package admin

import (
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"

	"github.com/swarnakumar/go-identity/db/sql/sqlc"
	"github.com/swarnakumar/go-identity/web/templates"
	"github.com/swarnakumar/go-identity/web/webapp"
)

var userChangeListingTemplate = templates.ParsePartial(
	"changes-listing.html",
	"admin/changes-listing.html",
	"components/pagination.html")

func RenderChanges(s webapp.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email, _ := url.QueryUnescape(s.URLParam(r, "email"))
		pageQuery := r.URL.Query().Get("page")

		var limit int32 = 10
		page := 1
		if pageQuery != "" {
			p, _ := strconv.ParseInt(pageQuery, 10, 32)
			page = int(p)
		}

		ctx := r.Context()

		c1 := make(chan []*sqlc.UserChange)
		c2 := make(chan int64)

		logger := s.GetLogger()
		db := s.GetDbClient()

		go func() {

			offset := (int32(page) - 1) * limit
			changes, err := db.Users.GetChangesForUser(ctx, email, offset, limit)
			if err != nil {
				logger.Errorw("Unable to get changes", "Error", err)
			}
			c1 <- changes
		}()

		go func() {
			count, err := db.Users.GetChangesCount(ctx, &email)
			if err != nil {
				logger.Errorw("Unable to get changes", "Error", err)
			}
			c2 <- count
		}()

		var listing []*sqlc.UserChange
		var count int64

		for i := 0; i < 2; i++ {
			select {
			case res := <-c1:
				listing = res
			case c := <-c2:
				count = c
			}
		}

		if listing == nil {
			if count == 0 {
				s.RenderErrorAlert(w, "No changes were EVER made!!!")
			} else {
				s.RenderErrorAlert(w, "Unable to Fetch Changes!!!")
			}
			return
		}

		pageCount := int(math.Ceil(float64(count) / float64(limit)))
		link := fmt.Sprintf("/admin/users/%s/changes", url.QueryEscape(email))
		params := map[string]interface{}{
			"Changes":        listing,
			"PageCount":      pageCount,
			"CurrPage":       page,
			"Link":           link,
			"PageContentDiv": "changes-listing-table"}

		s.ExecuteTemplate(w, r, userChangeListingTemplate, &params)
	}
}
