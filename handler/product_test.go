package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"order-system/utils"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateCaseFailed(t *testing.T) {
	tearDown()
	setup()
	var (
		reqJSON = `{"product":{"description":"product 3","price":300,"quantity":300}}`
	)
	req := httptest.NewRequest(http.MethodPost, "/api/products", strings.NewReader(reqJSON))
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", authHeader(utils.GenerateJWT(1)))
	h.Register(e)
	resp, err := e.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

//only vendor can create product
func TestCreateCaseSuccess(t *testing.T) {
	tearDown()
	setup()
	var (
		reqJSON = `{"product":{"description":"product 3","price":300,"quantity":300}}`
	)
	req := httptest.NewRequest(http.MethodPost, "/api/products", strings.NewReader(reqJSON))
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", authHeader(utils.GenerateJWT(2)))
	h.Register(e)
	resp, _ := e.Test(req, -1)
	if assert.Equal(t, http.StatusOK, resp.StatusCode) {
		body, _ := ioutil.ReadAll(resp.Body)
		m := responseMap(body, "product")
		assert.Equal(t, "product 3", m["description"])
	}
}

func TestGetProductCaseSuccess(t *testing.T) {
	tearDown()
	setup()
	url := "/api/products/" + strconv.Itoa(1)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Content-type", "application/json")
	h.Register(e)

	resp, err := e.Test(req, -1)
	assert.NoError(t, err)

	if assert.Equal(t, http.StatusOK, resp.StatusCode) {
		body, _ := ioutil.ReadAll(resp.Body)
		var a singleProductResponse
		err := json.Unmarshal(body, &a)
		m := responseMap(body, "product")
		assert.Equal(t, "product 1 description", m["description"])
		assert.Equal(t, float64(1), m["id"])
		assert.NoError(t, err)
	}
}

func TestGetProductCaseFailed(t *testing.T) {
	tearDown()
	setup()
	url := "/api/products/" + strconv.Itoa(5)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Content-type", "application/json")
	h.Register(e)

	resp, err := e.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestGetListProductCaseSuccess(t *testing.T) {
	tearDown()
	setup()
	url := "/api/products"
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", authHeader(utils.GenerateJWT(2)))
	h.Register(e)

	resp, err := e.Test(req, -1)
	assert.NoError(t, err)

	if assert.Equal(t, http.StatusOK, resp.StatusCode) {
		body, _ := ioutil.ReadAll(resp.Body)
		var a listProductResponse
		err := json.Unmarshal(body, &a)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), a.Count)
	}
}
