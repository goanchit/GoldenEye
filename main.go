package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"goldeneye.com/m/v2/api"
	"goldeneye.com/m/v2/config"
	"goldeneye.com/m/v2/consumers"
)

func main() {
	r := gin.Default()

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading env file")
	}

	mongoClient := config.ConnectDb()

	defer func() {
		if err := mongoClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	consumer, err := consumers.NewConsumer("AUTHOR_POST")
	consumer2, err := consumers.NewConsumer("AUTHOR_STATUS_JOB")

	defer consumer.Close()
	defer consumer2.Close()

	if err != nil {
		log.Fatal(err)
	}

	// Define multiple Consumers to distribute the work across different works

	// Consumer for author post
	go consumer.AuthorPostConsumer()

	// Consumer for author daily subscription job
	go consumer2.AuthorUpdateSubscription()

	api.RouteHander(r)

	if err := r.Run(":8000"); err != nil {
		log.Fatalln("Failed to run service")
	}
}