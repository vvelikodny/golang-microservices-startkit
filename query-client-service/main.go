package main

import (
	"github.com/nats-io/go-nats"
	"github.com/vvelikodny/golang-microservices-test/query-client-service/config"
	"log"
)

func main() {
	nc, err := nats.Connect(config.NatsURL)
	if err != nil {
		panic(err)
	}
	log.Printf("Connected to %s", nc.ConnectedUrl())

	app := NewApp(nc)
	app.Run()
}
