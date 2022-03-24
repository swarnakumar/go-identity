package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/swarnakumar/go-identity/db/sql/sqlc"
	"github.com/swarnakumar/go-identity/web/server"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleMe(t *testing.T) {
	ctx := context.Background()

	s := server.New(ctx)
	defer s.Close()

	s.InitRouter()

	apiRouter := MakeRouter(s)
	s.GetRouter().Mount("/api", apiRouter)

	r := s.GetRouter()

	// No token. Should return a 401
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/me", nil)
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	// Now token
	claims := map[string]interface{}{"user": email}
	token := s.GetJWT().Create(claims)
	rr = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/me", nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	r.ServeHTTP(rr, req)

	// Good token. Everything must be nice
	assert.Equal(t, http.StatusOK, rr.Code)
	var jsonResp1 sqlc.GetUserByEmailRow
	err := json.NewDecoder(rr.Body).Decode(&jsonResp1)
	assert.Nil(t, err)
	assert.Equal(t, email, jsonResp1.Email)

	// Random token
	rr = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/me", nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", "sdjkvndkjvndkfjv "))
	r.ServeHTTP(rr, req)

	// 401
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}
