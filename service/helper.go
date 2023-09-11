package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"

	"goldeneye.com/m/v2/common"
)

// Common Function To Publish To Q
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
