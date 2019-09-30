package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database

// Init ...
func Init(ctx context.Context, dbName, dbURI string) (err error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	client, err := mongo.NewClient(options.Client().ApplyURI(dbURI))
	if err != nil {
		return
	}

	if err = client.Connect(ctx); err != nil {
		return
	}

	db = client.Database(dbName)
	return
}

// FindOne ...
func FindOne(ctx context.Context, doc Document) error {
	return db.Collection(doc.Collection()).FindOne(ctx, doc.Where().bson()).Decode(doc)
}

// Find ...
func Find(ctx context.Context, doc Document) (*mongo.Cursor, error) {
	return db.Collection(doc.Collection()).Find(ctx, doc.Where().bson())
}

// InsertOne ...
func InsertOne(ctx context.Context, doc Document) (id primitive.ObjectID, err error) {
	result, err := db.Collection(doc.Collection()).InsertOne(ctx, doc)
	if err != nil {
		return
	}

	id = result.InsertedID.(primitive.ObjectID)
	return
}
