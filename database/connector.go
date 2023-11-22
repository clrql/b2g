package database

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database

func Conn() (Database *mongo.Database, Error error) {
	if db == nil {
		client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv("MONGO_URI")))
		if err != nil {
			return nil, err
		}
		db = client.Database(os.Getenv("MONGO_DATABASE"))
	}
	return db, nil
}
