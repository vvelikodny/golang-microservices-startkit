package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nats-io/go-nats"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/golang/protobuf/proto"
	"github.com/vvelikodny/golang-microservices-test/query-client-service/config"
	"github.com/vvelikodny/golang-microservices-test/query-client-service/errors"
	"github.com/vvelikodny/golang-microservices-test/query-client-service/news"
)

type App struct {
	Queue  *nats.Conn
	Router *mux.Router
}

func (app *App) Init(natsURL string) {
	nc, err := nats.Connect(natsURL)
	if err != nil {
		panic(err)
	}
	log.Printf("Connected to %s", nc.ConnectedUrl())
	app.Queue = nc

	app.Router = mux.NewRouter()
	app.initializeRoutes()
}

func (app *App) initializeRoutes() {
	app.Router.HandleFunc("/news", app.CreateNewsHandler).Methods(http.MethodPost)
	app.Router.HandleFunc("/news/{id}", app.GetNewsByIdHandler).Methods(http.MethodGet)
}

func (app *App) Run() {
	log.Fatal(http.ListenAndServe(":8080", app.Router))
}

func (app *App) CreateNewsHandler(w http.ResponseWriter, r *http.Request) {
	var n news.News
	if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
		errors.HttpError(w, fmt.Sprintf("JSON parsing error: %v", err), http.StatusInternalServerError)
		return
	}

	if _, err := govalidator.ValidateStruct(n); err != nil {
		errors.HttpError(w, fmt.Sprintf("New entity validation error: %v", err), http.StatusInternalServerError)
		return
	}

	requestMsg, err := proto.Marshal(&n)
	if err != nil {
		errors.HttpError(w, fmt.Sprintf("Entity marshal error: %v", err), http.StatusInternalServerError)
		return
	}

	msg, err := app.Queue.Request(config.CreateNewsChannel, requestMsg, 500*time.Millisecond)
	if err != nil {
		errors.HttpError(w, fmt.Sprintf("request error: %v", err), http.StatusInternalServerError)
		return
	}

	if err := proto.Unmarshal(msg.Data, &n); err != nil {
		errors.HttpError(w, fmt.Sprintf("Unmarshal error: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(n); err != nil {
		errors.HttpError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app *App) GetNewsByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	id, err := strconv.ParseInt(key, 10, 64)
	if err != nil {
		errors.HttpError(w, fmt.Sprintf("request error: %v", err), http.StatusInternalServerError)
		return
	}

	requestMsg, err := proto.Marshal(&news.GetNewsByIdRequest{Id: id})
	if err != nil {
		errors.HttpError(w, fmt.Sprintf("Entity marshal error: %v", err), http.StatusInternalServerError)
		return
	}

	msg, err := app.Queue.Request(config.GetNewsChannel, requestMsg, 10*time.Second)
	if err != nil {
		errors.HttpError(w, fmt.Sprintf("request error: %v", err), http.StatusInternalServerError)
		return
	}

	var n news.News
	if err := proto.Unmarshal(msg.Data, &n); err != nil {
		errors.HttpError(w, fmt.Sprintf("Unmarshal error: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(n); err != nil {
		errors.HttpError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
