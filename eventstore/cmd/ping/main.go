// Program used to test connection with DB
package main

import (
	"context"

	"eventstore/pkg/log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func Ping(uri string) {
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	// Ping the primary
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}

	log.PrintLn("transport", "mongodb", "msg", "succesfully pinged and connected to the DB")
}

func main() {
	uri := "mongodb://admin:admin123@eventstore-db-service:27017"
	Ping(uri)
}
