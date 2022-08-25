// Holds EventRepository definition and it's CRUD implementation
package repository

import (
	"context"
	"eventstore/pkg/grpc/pb"
	"eventstore/pkg/log"

	kitlog "github.com/go-kit/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EventRepository struct {
	db     *mongo.Database
	logger kitlog.Logger
}

func MakeEventRepository() EventRepository {
	uri := "mongodb://admin:admin123@eventstore-db-service:27017"
	reg := bson.NewRegistryBuilder().Build()
	client_opts := options.Client().ApplyURI(uri).SetRegistry(reg)

	client, err := mongo.Connect(context.Background(), client_opts)

	if err != nil {
		panic(err)
	}

	db := client.Database("eventstore")

	return EventRepository{
		db:     db,
		logger: log.MakeLogger(),
	}
}

func (repo EventRepository) Create(ctx context.Context, event *pb.Event) error {
	doc := bson.D{
		{"event_type", event.GetEventType()},
		{"aggregate_id", event.GetAggregateId()},
		{"aggregate_type", event.GetAggregateType()},
		{"event_data", event.GetEventData()},
		{"channel_name", event.GetChannelName()},
	}
	_, err := repo.db.Collection("events").InsertOne(ctx, doc)

	if err != nil {
		return err
	}
	return nil
}

func (repo EventRepository) Get(ctx context.Context, id string) (*pb.Event, error) {
	var results *pb.Event
	err := repo.db.Collection("events").FindOne(ctx, bson.M{"_id": id}).Decode(&results)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (repo EventRepository) Index(ctx context.Context) ([]*pb.Event, error) {
	cursor, err := repo.db.Collection("events").Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var results []*pb.Event
	// check for errors in the conversion
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}
