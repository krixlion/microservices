// Holds EventRepository definition and it's CRUD implementation
package repository

import (
	"context"
	"eventstore/pkg/grpc/pb"
	"eventstore/pkg/log"

	kitlog "github.com/go-kit/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EventRepository struct {
	db     *mongo.Database
	logger kitlog.Logger
}

func MakeEventRepository() EventRepository {
	uri := "mongodb://admin:admin123@eventstore-db-service:27017"

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	db := client.Database("change-me")

	return EventRepository{
		db:     db,
		logger: log.MakeLogger(),
	}
}

func (s EventRepository) Create(data *pb.Event) error {
	return nil
}

func (s EventRepository) Get(id string) (*pb.Event, error) {
	return nil, nil
}

func (s EventRepository) Index() ([]*pb.Event, error) {
	return nil, nil
}
