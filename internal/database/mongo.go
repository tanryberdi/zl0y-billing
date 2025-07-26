package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func NewMongoDB(uri, dbName string) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, _ := mongo.Connect(options.Client().ApplyURI(uri))
	//defer func() {
	//	if err := client.Disconnect(ctx); err != nil {
	//		panic(err)
	//	}
	//}()

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	db := client.Database(dbName)
	mongoDB := &MongoDB{
		Client:   client,
		Database: db,
	}

	// Create indexes
	if err := mongoDB.createIndexes(); err != nil {
		return nil, fmt.Errorf("failed to create indexes: %w", err)
	}

	return mongoDB, nil
}

func (m *MongoDB) Disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return m.Client.Disconnect(ctx)
}

func (m *MongoDB) createIndexes() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := m.Database.Collection("reports")

	// Create index on report_id (unique)
	_, err := collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "report_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return fmt.Errorf("failed to create index on report_id: %w", err)
	}

	// Create index on user_id for faster queries - using bson.D for consistency
	_, err = collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "user_id", Value: 1}},
	})
	if err != nil {
		return fmt.Errorf("failed to create index on user_id: %w", err)
	}

	// Create index on client_generated_id for an anonymous report linking - using bson.D
	_, err = collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "client_generated_id", Value: 1}},
	})
	if err != nil {
		return fmt.Errorf("failed to create index on client_generated_id: %w", err)
	}

	// Create a compound index for efficient queries - using bson.D (ORDER MATTERS)
	_, err = collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "user_id", Value: 1},
			{Key: "created_at", Value: -1},
		},
	})

	return err
}
