package main

import (
	"github.com/vvelikodny/golang-microservices-test/query-client-service/config"
)

func main() {
	app := App{}
	app.Init(config.NatsUrl)
	app.Run()
}
