package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"

	"goldeneye.com/m/v2/common"
	"goldeneye.com/m/v2/database"
	"goldeneye.com/m/v2/models"
)

func PublishToQueue(ctx context.Context, queueName string, data interface{}) {

	bytes, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	messageString := fmt.Sprintf("%s", base64.StdEncoding.EncodeToString(bytes))

	success, errorString := common.SendMessageToQ(ctx, queueName, messageString)
	if !success {
		log.Fatalf("Error occurred while pushing to Queue %s : %s", queueName, errorString)
	}
	return
}

// Convert Array to Slices of size chunkSize
func chunkBy[T any](items []T, chunkSize int) (chunks [][]T) {
	for chunkSize < len(items) {
		items, chunks = items[chunkSize:], append(chunks, items[0:chunkSize:chunkSize])
	}
	return append(chunks, items)
}

func UpdateAuthorStatusJob(ctx context.Context) {
	client, err := database.NewClient()

	if err != nil {
		log.Fatalf("Failed to get database client %s", err)
		return
	}

	allAuthors, settings := client.GetAllAuthorData(context.TODO())

	// Create author list batches and pass to worker
	batches := chunkBy(allAuthors, 4)

	for _, v := range batches {
		m := make(map[string]interface{})

		m["data"] = v
		m["settings"] = settings

		PublishToQueue(ctx, "AUTHOR_STATUS_JOB", m)
	}

	return
}

func UpdateGlobalSettings(ctx context.Context, data models.Settings) {
	client, err := database.NewClient()
	if err != nil {
		log.Fatalf("Failed to get database client %s", err)
		return
	}

	client.UpdateGlobalSettings(ctx, data)
	return
}
