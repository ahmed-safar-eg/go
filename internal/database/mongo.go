package database

import (
	"context"
	"project/internal/config"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var MongoDB *mongo.Database

func InitializeMongo(cfg config.MongoConfig) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.URI))
    if err != nil {
        return err
    }
    MongoClient = client
    MongoDB = client.Database(cfg.DBName)
    return nil
}

func GetMongoDB() *mongo.Database {
    return MongoDB
}

func CloseMongo() error {
    if MongoClient != nil {
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()
        return MongoClient.Disconnect(ctx)
    }
    return nil
}

// EnsureIndexes creates required indexes (e.g., unique email on users)
func EnsureIndexes() error {
    if MongoDB == nil {
        return nil
    }
    users := MongoDB.Collection("users")
    // unique index on email
    // Partial unique index: only for documents without deleted_at
    idx := mongo.IndexModel{
        Keys: bson.D{{Key: "email", Value: 1}},
        Options: options.Index().
            SetUnique(true).
            SetName("uniq_email").
            SetPartialFilterExpression(bson.D{{Key: "deleted_at", Value: bson.D{{Key: "$eq", Value: nil}}}}),
    }
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    _, err := users.Indexes().CreateOne(ctx, idx)
    return err
}


