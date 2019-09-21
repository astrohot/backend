package database

import (
	"context"
	"time"

	"github.com/astrohot/backend/internal/model/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Database config
const (
	DBName = "astrohot"
	URI    = "mongodb+srv://admin:admin@cluster0-srljh.mongodb.net/test?retryWrites=true&w=majority"
)

// Create ...
func Create(ctx context.Context) (db *DB, err error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	c, err := mongo.NewClient(options.Client().ApplyURI(URI))
	if err != nil {
		return
	}

	c.Connect(ctx)
	db = &DB{client: c}
	return
}

// DB ...
type DB struct {
	client *mongo.Client
}

// Collection ...
func (db *DB) Collection(coll string) *mongo.Collection {
	return db.client.Database(DBName).Collection(coll)
}

// Exists ...
func (db *DB) Exists(ctx context.Context, coll *mongo.Collection, where interface{}) (exists bool) {
	if err := coll.FindOne(ctx, where).Err(); err == nil {
		exists = true
	}

	return
}

// GetUserByEmail ...
func (db *DB) GetUserByEmail(ctx context.Context, email string) (user.User, error) {
	var u user.User
	coll := db.Collection("users")
	where := bson.M{"email": email}

	err := coll.FindOne(ctx, where).Decode(&u)
	return u, err
}
