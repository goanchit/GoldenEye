package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func getMongoConnectionString() string {
	return fmt.Sprintf("mongodb+srv://%s:%s@cluster0.9jgzm.mongodb.net/?retryWrites=true&w=majority", os.Getenv("MONGO_USER"), os.Getenv("MONGO_PASSWORD"))
}

func ConnectDb() *mongo.Client {
	connectionString := getMongoConnectionString()

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(connectionString).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}

	log.Println("Database Connected")

	Client = client
	return client
}
