package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/swarnakumar/go-identity/db"
	"github.com/swarnakumar/go-identity/web/server"
	"net/http"
	"net/http/httptest"
	"testing"
)

var email = "abc-refresh@xyz.com"
var pwd = "qwerty1234@96219authmiddleware"

type resp struct {
	Token string `json:"token"`
}

func TestMain(m *testing.M) {
	ctx := context.Background()

	dbClient := db.New(ctx)
	defer dbClient.Close()

	_, _ = dbClient.Users.Create(ctx, email, pwd, false, nil)
	m.Run()
}

func TestGenerateToken(t *testing.T) {
	ctx := context.Background()

	s := server.New(ctx)
	defer s.Close()

	s.InitRouter()

	apiRouter := MakeRouter(s)
	s.GetRouter().Mount("/api", apiRouter)

	rr := httptest.NewRecorder()
	handler := HandleGenerateToken(s)

	email := "abc@xyz.com"
	pwd := "qwerty1234@96219validtoken"

	body := getTokenParams{Email: email, Password: pwd}

	reqBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/api/generate-token", bytes.NewBuffer(reqBody))
	handler.ServeHTTP(rr, req)

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

func TestHandleRefreshToken(t *testing.T) {
	ctx := context.Background()

	s := server.New(ctx)
	defer s.Close()

	s.InitRouter()
	apiRouter := MakeRouter(s)
	s.GetRouter().Mount("/api", apiRouter)

	r := s.GetRouter()
	//r.Post("/refresh", HandleRefreshToken(s))

	claims := map[string]interface{}{"user": email}
	token := s.GetJWT().Create(claims)

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/refresh-token", nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	r.ServeHTTP(rr, req)

	// Good token. Everything must be nice
	assert.Equal(t, http.StatusCreated, rr.Code)
	var jsonResp1 resp
	json.NewDecoder(rr.Body).Decode(&jsonResp1)
	// Must be a non-empty token
	assert.NotEqual(t, jsonResp1.Token, "")

	// Now without a token.
	rr = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/refresh-token", nil)
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	var jsonResp2 resp
	json.NewDecoder(rr.Body).Decode(&jsonResp2)
	// Must be a non-empty token
	assert.Equal(t, jsonResp2.Token, "")
}

func TestHandleVerifyToken(t *testing.T) {
	ctx := context.Background()

	s := server.New(ctx)
	defer s.Close()

	s.InitRouter()
	apiRouter := MakeRouter(s)
	s.GetRouter().Mount("/api", apiRouter)

	r := s.GetRouter()

	type resp struct {
		Valid bool `json:"valid"`
	}

	claims := map[string]interface{}{"user": email}
	token := s.GetJWT().Create(claims)
	body := verifyTokenParams{Token: token}
	reqBody, _ := json.Marshal(body)

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/verify-token", bytes.NewBuffer(reqBody))
	r.ServeHTTP(rr, req)

	// Good token. Everything must be nice
	assert.Equal(t, http.StatusOK, rr.Code)
	var jsonResp1 resp
	json.NewDecoder(rr.Body).Decode(&jsonResp1)
	// This is a good token. Valid MUST be true!!!
	assert.True(t, jsonResp1.Valid)

	// Bad token.
	body = verifyTokenParams{Token: "djcnejfbvwkn"}
	reqBody, _ = json.Marshal(body)

	rr = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/verify-token", bytes.NewBuffer(reqBody))
	r.ServeHTTP(rr, req)
	// Token's passed. So 200.
	assert.Equal(t, http.StatusOK, rr.Code)
	var jsonResp2 resp
	json.NewDecoder(rr.Body).Decode(&jsonResp2)
	// This is a bad token. Valid MUST be false!!!
	assert.False(t, jsonResp2.Valid)

	// No TOKEN
	rr = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/verify-token", nil)
	r.ServeHTTP(rr, req)
	// Token's not passed. So 400.
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Random Body, i.e., without TOKEN
	email := "abc@xyz.com"
	pwd := "qwerty1234@96219validtoken"

	randomBody := getTokenParams{Email: email, Password: pwd}
	reqBody, _ = json.Marshal(randomBody)

	rr = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/verify-token", bytes.NewBuffer(reqBody))
	r.ServeHTTP(rr, req)
	// Again 400.
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
