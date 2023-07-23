package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type Database struct {
	Client     *mongo.Client
	Collection *mongo.Collection
}

func NewDatabase(mongoURI, dbName, collectionName string) (*Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, fmt.Errorf("error connecting to MongoDB: %v", err)
	}

	log.Println("Connected to MongoDB!")

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("error pinging MongoDB: %v", err)
	}

	db := client.Database(dbName)

	collections, err := db.ListCollectionNames(ctx, bson.M{"name": collectionName})
	if err != nil {
		return nil, fmt.Errorf("error listing collections: %v", err)
	}

	if len(collections) == 0 {
		err = db.CreateCollection(ctx, collectionName)
		if err != nil {
			return nil, fmt.Errorf("error creating collection: %v", err)
		}
	}

	database := &Database{
		Client:     client,
		Collection: db.Collection(collectionName),
	}

	return database, nil
}

func (db *Database) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db.Client.Disconnect(ctx)
}
