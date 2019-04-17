package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/nats-io/go-nats"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/vvelikodny/golang-microservices-test/query-client-service/news"
)

type MockNatsConnection struct {
	mock.Mock
}

func (mock *MockNatsConnection) Request(subj string, data []byte, timeout time.Duration) (*nats.Msg, error) {
	args := mock.Called(subj, data, timeout)
	return args.Get(0).(*nats.Msg), args.Error(1)
}

func TestCreateNewsError(t *testing.T) {
	nc := &MockNatsConnection{}
	nc.On("Request", mock.Anything, mock.Anything, mock.Anything).Return(&nats.Msg{}, errors.New("error fro test"))

	app := NewApp(nc)

	req, _ := http.NewRequest(http.MethodPost, "/news", bytes.NewBuffer([]byte(`{"Title": "Hello, Golang!"}`)))
	response := executeRequest(t, app, req)

	checkResponseCode(t, http.StatusInternalServerError, response.Code)

	var m map[string]interface{}
	require.Error(t, json.Unmarshal(response.Body.Bytes(), &m))
}

func TestCreateNews(t *testing.T) {
	nc := &MockNatsConnection{}

	data, err := proto.Marshal(&news.News{Id: 1, Title: "title"})
	require.NoError(t, err)
	nc.On("Request", mock.Anything, mock.Anything, mock.Anything).Return(&nats.Msg{Data: data}, nil)

	app := NewApp(nc)

	req, _ := http.NewRequest(http.MethodPost, "/news", bytes.NewBuffer([]byte(`{"Title": "Hello, Golang!"}`)))
	response := executeRequest(t, app, req)

	t.Logf("%+v", string(response.Body.Bytes()))
	require.Equal(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &m))

	require.Contains(t, m, "id")
	require.True(t, m["id"].(float64) > 0)
}

func TestGetNews(t *testing.T) {
	data, err := proto.Marshal(&news.News{Id: 1, Title: "title"})
	require.NoError(t, err)

	nc := &MockNatsConnection{}
	nc.On("Request", mock.Anything, mock.Anything, mock.Anything).Return(&nats.Msg{Data: data}, nil)

	app := NewApp(nc)

	req, _ := http.NewRequest(http.MethodGet, "/news/1", nil)
	response := executeRequest(t, app, req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m news.News
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &m))
	require.Equal(t, m.Id, int64(1))
	require.Equal(t, m.Title, "title")
}

func executeRequest(t *testing.T, app *App, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d, but got %d\n", expected, actual)
	}
}
