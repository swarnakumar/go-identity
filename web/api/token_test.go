package api

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/swarnakumar/go-identity/web/server"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGenerateToken(t *testing.T) {
	ctx := context.Background()

	s := server.New(ctx)
	defer s.Close()

	s.InitRouter()

	apiRouter := MakeApiRouter(s)
	s.GetRouter().Mount("/api", apiRouter)

	rr := httptest.NewRecorder()
	handler := HandleGenerateToken(s)

	email := "abc@xyz.com"
	pwd := "qwerty1234@96219validtoken"

	body := getTokenParams{Email: email, Password: pwd}

	reqBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/api/generate-token", bytes.NewBuffer(reqBody))
	handler.ServeHTTP(rr, req)

	type resp struct {
		Token string `json:"token"`
	}

	// User doesnt exist. Must return a 401.
	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	var jsonResp1 resp
	json.NewDecoder(rr.Body).Decode(&jsonResp1)
	// Must be an empty token
	assert.Equal(t, jsonResp1.Token, "")

	// Creating a new user
	s.GetDbClient().Users.Create(ctx, email, pwd, false, nil)
	rr = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/generate-token", bytes.NewBuffer(reqBody))
	handler.ServeHTTP(rr, req)

	// Now the user exists!!!
	assert.Equal(t, http.StatusCreated, rr.Code)

	var jsonResp2 resp
	json.NewDecoder(rr.Body).Decode(&jsonResp2)
	// Must be a non-empty token
	assert.NotEqual(t, jsonResp2.Token, "")

	// No password
	body = getTokenParams{Email: email}

	reqBody, _ = json.Marshal(body)
	rr = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/generate-token", bytes.NewBuffer(reqBody))
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Wrong password
	body = getTokenParams{Email: email, Password: "djksnvkrjeigbvierjbv"}

	reqBody, _ = json.Marshal(body)
	rr = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/generate-token", bytes.NewBuffer(reqBody))
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}
