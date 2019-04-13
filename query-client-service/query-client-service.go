package main

import (
	"github.com/gorilla/mux"
	"github.com/nats-io/go-nats"
	"github.com/vvelikodny/golang-microservices-test/config"
	"github.com/vvelikodny/golang-microservices-test/news"
	"log"
	"net/http"
)

func main() {
	nc, err := nats.Connect(config.NatsUrl)
	if err != nil {
		panic(err)
	}

	defer nc.Close()
	log.Printf("Connected to %s", nc.ConnectedUrl())

	r := mux.NewRouter()

	r.HandleFunc("/news", news.CreateNewsHandler(nc)).Methods(http.MethodPost)
	r.HandleFunc("/news/{id}", news.GetNewsByIdHandler(nc)).Methods(http.MethodGet)

	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}
}
