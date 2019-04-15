package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vvelikodny/golang-microservices-test/news"
)

var app App

func TestMain(m *testing.M) {
	app = App{}
	app.Init("")
	os.Exit(m.Run())
}

func TestCreateNews(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPost, "/news", bytes.NewBuffer([]byte(`{"Title": "Hello, Golang!"}`)))
	response := executeRequest(t, req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &m))

	require.Contains(t, m, "id")
	require.True(t, m["id"].(float64) > 0)
}

func TestGetNews(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/news/1", nil)
	response := executeRequest(t, req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m news.News
	assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &m))
	assert.Equal(t, m.Id, int64(1))
	assert.Equal(t, m.Title, "Hello, Golang!")
}

func executeRequest(t *testing.T, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d, but got %d\n", expected, actual)
	}
}
