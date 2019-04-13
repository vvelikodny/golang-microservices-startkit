package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/nats-io/go-nats"
	"github.com/vvelikodny/golang-microservices-test/storage-service/news"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// TODO to .env
const PostgresUrl = "postgresql://postgres:postgres@db:5432/postgres?sslmode=disable"
const NatsUrl = "nats://natsd:4222"

const CreateNewsChannel = "news-create"
const GetNewsChannel = "get-news"

var nc *nats.Conn

func main() {
	db, err := gorm.Open("postgres", PostgresUrl)
	if err != nil {
		log.Fatal(fmt.Sprintf("err: %v", err))
	}
	defer db.Close()

	// Automatically create the "news" table based on the Account model.
	db.AutoMigrate(&news.News{})

	nc, err = nats.Connect(NatsUrl)
	if err != nil {
		panic(err)
	}
	defer nc.Close()
	log.Printf("Connected to %s", nc.ConnectedUrl())

	nc.Subscribe(CreateNewsChannel, func(m *nats.Msg) {
		var n news.News
		proto.Unmarshal(m.Data, &n)
		if err = db.Create(&n).Error; err != nil {
			panic(err)
		}
		response, _ := proto.Marshal(&n)
		nc.Publish(m.Reply, response)
	})

	nc.Subscribe(GetNewsChannel, func(m *nats.Msg) {
		var id news.GetNewsByIdRequest
		proto.Unmarshal(m.Data, &id)

		var n news.News
		if err = db.First(&n, id.Id).Error; err != nil {
			// TODO handle errors through NATS message
			panic(err)
		}
		response, _ := proto.Marshal(&n)
		nc.Publish(m.Reply, response)
	})

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGINT)

	<-c
}
