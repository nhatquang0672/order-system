package handler

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"order-system/utils"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignUpCaseSuccess(t *testing.T) {
	tearDown()
	setup()
	var (
		reqJSON = `{"user":{"role":"user","email":"alice@realworld.io","password":"secret"}}`
	)
	req := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(reqJSON))
	req.Header.Set("Content-type", "application/json")
	h.Register(e)
	resp, _ := e.Test(req, -1)
	if assert.Equal(t, http.StatusCreated, resp.StatusCode) {
		body, _ := ioutil.ReadAll(resp.Body)
		m := responseMap(body, "user")
		assert.Equal(t, "user", m["role"])
		assert.Equal(t, "alice@realworld.io", m["email"])
		assert.NotEmpty(t, m["token"])
	}
}

func TestLoginCaseSuccess(t *testing.T) {
	tearDown()
	setup()
	var (
		reqJSON = `{"user":{"email":"user1@realworld.io","password":"secret"}}`
	)
	req := httptest.NewRequest(http.MethodPost, "/api/users/login", strings.NewReader(reqJSON))
	req.Header.Set("Content-type", "application/json")
	h.Register(e)
	resp, _ := e.Test(req, -1)
	if assert.Equal(t, http.StatusOK, resp.StatusCode) {
		body, _ := ioutil.ReadAll(resp.Body)
		m := responseMap(body, "user")
		assert.Equal(t, "user1", m["username"])
		assert.Equal(t, "user1@realworld.io", m["email"])
	}
}

func TestLoginCaseFailed(t *testing.T) {
	tearDown()
	setup()
	var (
		reqJSON = `{"user":{"email":"userx@realworld.io","password":"secret"}}`
	)
	req := httptest.NewRequest(http.MethodPost, "/api/users/login", strings.NewReader(reqJSON))
	req.Header.Set("Content-type", "application/json")
	h.Register(e)
	resp, err := e.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
}
func TestCurrentUserCaseSuccess(t *testing.T) {
	tearDown()
	setup()
	req := httptest.NewRequest(http.MethodGet, "/api/user", nil)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", authHeader(utils.GenerateJWT(1)))
	h.Register(e)
	resp, err := e.Test(req, -1)
	assert.NoError(t, err)
	if assert.Equal(t, http.StatusOK, resp.StatusCode) {
		body, _ := ioutil.ReadAll(resp.Body)
		m := responseMap(body, "user")
		assert.Equal(t, "user1", m["username"])
		assert.Equal(t, "user1@realworld.io", m["email"])
		assert.NotEmpty(t, m["token"])
	}
}

func TestCurrentUserCaseInvalid(t *testing.T) {
	tearDown()
	setup()
	req := httptest.NewRequest(http.MethodGet, "/api/user", nil)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", authHeader(utils.GenerateJWT(100)))
	h.Register(e)
	resp, err := e.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestUpdateUserEmail(t *testing.T) {
	tearDown()
	setup()
	var (
		user1UpdateReq = `{"user":{"email":"user1@user1.me"}}`
	)
	req := httptest.NewRequest(http.MethodPut, "/api/user", strings.NewReader(user1UpdateReq))
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", authHeader(utils.GenerateJWT(1)))
	h.Register(e)
	resp, err := e.Test(req, -1)
	assert.NoError(t, err)

	if assert.Equal(t, http.StatusOK, resp.StatusCode) {
		body, _ := ioutil.ReadAll(resp.Body)
		m := responseMap(body, "user")
		assert.Equal(t, "user1", m["username"])
		assert.Equal(t, "user1@user1.me", m["email"])
		assert.NotEmpty(t, m["token"])
	}
}

func TestUpdateUserMultipleField(t *testing.T) {
	tearDown()
	setup()
	var (
		user1UpdateReq = `{"user":{"username":"user11", "email":"user11@user11.me"}}`
	)
	req := httptest.NewRequest(http.MethodPut, "/api/user", strings.NewReader(user1UpdateReq))
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", authHeader(utils.GenerateJWT(1)))
	h.Register(e)
	resp, err := e.Test(req, -1)
	assert.NoError(t, err)

	if assert.Equal(t, http.StatusOK, resp.StatusCode) {
		body, _ := ioutil.ReadAll(resp.Body)
		m := responseMap(body, "user")
		assert.Equal(t, "user11", m["username"])
		assert.Equal(t, "user11@user11.me", m["email"])
		assert.NotEmpty(t, m["token"])
	}
}
