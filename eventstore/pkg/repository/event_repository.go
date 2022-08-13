// Holds EventRepository definition and it's CRUD implementation
package repository

import (
	"context"
	"eventstore/pkg/grpc/pb"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EventRepository struct {
	db *mongo.Database
}

func MakeEventRepository() EventRepository {
	uri := "mongodb://admin:admin123@eventstore-db-service:27017"

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	db := client.Database("")

	return EventRepository{
		db: db,
	}
}

func (s EventRepository) Create(data *pb.Event) error {
	return nil
}

func (s EventRepository) Get(id string) (*pb.Event, error) {
	return nil, nil
}
