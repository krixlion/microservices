package event_repository

import (
	"context"
	"fmt"
	"os"

	"eventstore/internal/pb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbName = "eventstore"
)

type EventRepository struct {
	db *mongo.Database
	// logger kitlog.Logger
}

// MakeEventRepository takes connection data from environment variables and forms a uri which it later uses for mongodb connection
// Panics if it fails to connect
func MakeEventRepository() EventRepository {
	pass := os.Getenv("DB_PASS")
	port := os.Getenv("DB_PORT")
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")

	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s", user, pass, host, port)

	reg := bson.NewRegistryBuilder().Build()
	client_opts := options.Client().ApplyURI(uri).SetRegistry(reg)

	client, err := mongo.Connect(context.Background(), client_opts)

	if err != nil {
		err = fmt.Errorf("failed to connect to the DB: %+w", err)
		panic(err)
	}

	db := client.Database(dbName)

	return EventRepository{
		db: db,
		// logger: log.MakeLogger(),
	}
}

// Close disconnects the client form MongoDB connection pool
func (repo EventRepository) Close(ctx context.Context) error {
	return repo.db.Client().Disconnect(ctx)
}

// Create inserts one event into DB
func (repo EventRepository) Create(ctx context.Context, event *pb.Event) (err error) {
	doc := bson.D{
		{"event_type", event.GetEventType()},
		{"aggregate_id", event.GetAggregateId()},
		{"aggregate_type", event.GetAggregateType()},
		{"event_data", event.GetEventData()},
		// {"channel_name", event.GetChannelName()},
	}
	_, err = repo.db.Collection("events").InsertOne(ctx, doc)

	if err != nil {
		return err
	}
	return nil
}

// Get returns the event with given ID
func (repo EventRepository) Get(ctx context.Context, id string) (event *pb.Event, err error) {
	err = repo.db.Collection("events").FindOne(ctx, bson.M{"_id": id}).Decode(&event)
	if err != nil {
		return nil, err
	}

	return event, nil
}

// Index returns all events
func (repo EventRepository) Index(ctx context.Context) (events []*pb.Event, err error) {
	cursor, err := repo.db.Collection("events").Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	// check for errors in the conversion
	if err = cursor.All(ctx, &events); err != nil {
		return nil, err
	}
	return events, nil
}
