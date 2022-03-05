package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"order-system/utils"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateCartItemCaseFailed(t *testing.T) {
	tearDown()
	setup()
	var (
		reqJSON = `{"cartItem":{"productid":5,"quantity":5}}`
	)
	req := httptest.NewRequest(http.MethodPost, "/api/carts", strings.NewReader(reqJSON))
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", authHeader(utils.GenerateJWT(1)))
	h.Register(e)
	resp, err := e.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestCreateCartItemCaseSuccess(t *testing.T) {
	tearDown()
	setup()
	var (
		reqJSON = `{"cartItem":{"productid":2,"quantity":5}}`
	)
	req := httptest.NewRequest(http.MethodPost, "/api/carts", strings.NewReader(reqJSON))
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", authHeader(utils.GenerateJWT(1)))
	h.Register(e)
	resp, err := e.Test(req, -1)
	assert.NoError(t, err)
	if assert.Equal(t, http.StatusOK, resp.StatusCode) {
		body, _ := ioutil.ReadAll(resp.Body)
		var a singleCartItemResponse
		err := json.Unmarshal(body, &a)
		m := responseMap(body, "cartItem")
		assert.Equal(t, float64(5), m["quantity"])
		assert.NoError(t, err)
	}
}

func TestGetCartItemsCaseSuccess(t *testing.T) {
	tearDown()
	setup()
	url := "/api/carts"
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", authHeader(utils.GenerateJWT(1)))
	h.Register(e)

	resp, err := e.Test(req, -1)
	assert.NoError(t, err)

	if assert.Equal(t, http.StatusOK, resp.StatusCode) {
		body, _ := ioutil.ReadAll(resp.Body)
		var a listCartItemsResponse
		err := json.Unmarshal(body, &a)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(a.CartItems))
		assert.Equal(t, uint(600), a.TotalPrice)
	}
}

func TestGetCartItemsCaseEmptySuccess(t *testing.T) {
	tearDown()
	setup()
	url := "/api/carts"
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", authHeader(utils.GenerateJWT(2)))
	h.Register(e)

	resp, err := e.Test(req, -1)
	assert.NoError(t, err)

	if assert.Equal(t, http.StatusOK, resp.StatusCode) {
		body, _ := ioutil.ReadAll(resp.Body)
		var a listCartItemsResponse
		err := json.Unmarshal(body, &a)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(a.CartItems))
	}
}

func TestGetCartItemCaseFailed(t *testing.T) {
	tearDown()
	setup()
	url := "/api/carts"
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", authHeader(utils.GenerateJWT(3)))
	h.Register(e)

	resp, err := e.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestPlaceOrderCaseFailed(t *testing.T) {
	tearDown()
	setup()
	url := "/api/orders"
	req := httptest.NewRequest(http.MethodPost, url, nil)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", authHeader(utils.GenerateJWT(2)))
	h.Register(e)

	resp, err := e.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestPlaceOrderCaseSuccess(t *testing.T) {
	tearDown()
	setup()
	url := "/api/orders"
	req := httptest.NewRequest(http.MethodPost, url, nil)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", authHeader(utils.GenerateJWT(1)))
	h.Register(e)

	resp, err := e.Test(req, -1)
	assert.NoError(t, err)

	if assert.Equal(t, http.StatusOK, resp.StatusCode) {
		body, _ := ioutil.ReadAll(resp.Body)
		var a singleOrderResponse
		err := json.Unmarshal(body, &a)
		assert.NoError(t, err)
		assert.Equal(t, "Paid", a.Order.Status)
		assert.Equal(t, uint(600), a.Order.TotalPrice)
		assert.Equal(t, 1, len(a.Order.Items))
	}
}

func TestGetOrderCaseSuccess(t *testing.T) {
	tearDown()
	setup()
	url := "/api/orders"
	req := httptest.NewRequest(http.MethodPost, url, nil)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", authHeader(utils.GenerateJWT(1)))
	h.Register(e)

	resp, err := e.Test(req, -1)
	assert.NoError(t, err)

	if assert.Equal(t, http.StatusOK, resp.StatusCode) {
		body, _ := ioutil.ReadAll(resp.Body)
		var a singleOrderResponse
		err := json.Unmarshal(body, &a)
		assert.NoError(t, err)
		assert.Equal(t, "Paid", a.Order.Status)
		assert.Equal(t, uint(600), a.Order.TotalPrice)
		assert.Equal(t, 1, len(a.Order.Items))
	}
}
