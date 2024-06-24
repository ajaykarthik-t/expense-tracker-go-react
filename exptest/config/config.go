// config/config.go

package config

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var Db *mongo.Database

const (
	DBName         = "ExpenseTracker"
	CollectionName = "Transactions"
)

func ConnectDB() error {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")

	// Connect to MongoDB
	var err error
	Client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return fmt.Errorf("error connecting to MongoDB: %v", err)
	}

	// Check the connection
	err = Client.Ping(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("error pinging MongoDB: %v", err)
	}

	Db = Client.Database(DBName)
	log.Println("Connected to MongoDB!")

	return nil
}
