package main

import (
	"github.com/vvelikodny/golang-microservices-test/config"
)

func main() {
	app := App{}
	app.Init(config.NatsUrl)
	app.Run()
}
