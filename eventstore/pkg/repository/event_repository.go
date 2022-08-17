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
	client_opts := options.Client().ApplyURI(uri)
	client_opts.Registry = reg

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

func (s EventRepository) Create(ctx context.Context, data *pb.Event) error {
	doc := bson.D{
		{"_id", data.EventId}, {"event_type", data.EventType}, {"aggregate_id", data.AggregateId},
	}
	_, err := s.db.Collection("events").InsertOne(ctx, doc)

	if err != nil {
		s.logger.Log("msg", "failed to create event", "db", "err", err)
		return err
	}

	return nil
}

func (s EventRepository) Get(ctx context.Context, id string) (*pb.Event, error) {
	cursor, err := s.db.Collection("events").Find(ctx, bson.D{{"_id": id}})
	if err != nil {
		s.logger.Log("msg", "failed to get event", "db", "err", err)
		return nil, err
	}

	var results []bson.M
	// check for errors in the conversion
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}
	return nil, nil
}

func (s EventRepository) Index(ctx context.Context) ([]*pb.Event, error) {
	cursor, err := s.db.Collection("events").Find(ctx, bson.D{})
	if err != nil {
		s.logger.Log("msg", "failed to get event", "db", "err", err)
		return nil, err
	}

	var results []bson.M
	// check for errors in the conversion
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}
	return nil, nil
}
