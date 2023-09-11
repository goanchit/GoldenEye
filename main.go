package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	ctx := context.Background()

	mongoClient := config.ConnectDb(ctx)

	defer func() {
		if err := mongoClient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	consumer, err := consumers.NewConsumer("AUTHOR_POST", mongoClient)
	consumer2, err := consumers.NewConsumer("AUTHOR_STATUS_JOB", mongoClient)

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

	api.RouteHander(r, mongoClient)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")
}
