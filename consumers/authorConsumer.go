package consumers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"goldeneye.com/m/v2/database"
	"goldeneye.com/m/v2/models"
)

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

func getMQConnectionString() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s/", os.Getenv("RABBITMQ_LOGIN"), os.Getenv("RABBITMQ_PASSWORD"), os.Getenv("RABBITMQ_HOST"), os.Getenv("RABBITMQ_PORT"))
}

func NewConsumer(queueName string) (*Consumer, error) {
	connection_string := getMQConnectionString()

	conn, err := amqp.Dial(connection_string)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	q, err := channel.QueueDeclare(
		queueName, true, false, false, false, nil,
	)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		conn:    conn,
		channel: channel,
		queue:   q,
	}, nil
}

func (c *Consumer) Close() {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *Consumer) AuthorPostConsumer() error {
	msgs, err := c.channel.Consume(
		c.queue.Name, "", true, false, false, false, nil,
	)
	if err != nil {
		return err
	}

	for d := range msgs {
		log.Printf("Received a new Message message: %s", d.Body)
		byteString, err := base64.StdEncoding.DecodeString(string(d.Body))

		if err != nil {
			log.Fatalf("Failed to decode base64 string: %s", err)
		}

		var messageBody models.MessageBody
		err = json.Unmarshal(byteString, &messageBody)
		if err != nil {
			log.Fatal(err)
		}
		// Upsert Post increment Count and fanout based on eligibility
		mongoClient, err := database.NewClient()
		isPremiumAuthor := mongoClient.InsertAuthorPost(context.TODO(), messageBody)

		if isPremiumAuthor {
			// Send Post to premium users
		} else {
			// Send Post to regular users
		}

		// Below code depicts number of random followers to be added
		randomFollowers := rand.Intn(10)

		mongoClient.UpdateAuthorFollowers(context.TODO(), messageBody.UserId, randomFollowers)

	}

	return nil
}

func (c *Consumer) AuthorUpdateSubscription() error {
	msgs, err := c.channel.Consume(
		c.queue.Name, "", false, false, false, false, nil,
	)
	if err != nil {
		return err
	}

	for d := range msgs {
		log.Printf("Received a new Message message: %s", d.Body)
		byteString, err := base64.StdEncoding.DecodeString(string(d.Body))

		if err != nil {
			log.Fatalf("Failed to decode base64 string: %s", err)
		}

		var messageBody models.AuthorPremiumJob
		err = json.Unmarshal(byteString, &messageBody)
		if err != nil {
			log.Fatal(err)
		}

		var premiumAuthorsList []string

		// For NonPremium Authors set premium flag to false
		var nonpremiumAuthorsList []string

		for _, v := range messageBody.Data {
			if messageBody.Settings.MinFollowers <= v.Followers {
				premiumAuthorsList = append(premiumAuthorsList, v.UUID)
			} else {
				nonpremiumAuthorsList = append(nonpremiumAuthorsList, v.UUID)
			}

		}

		mongoClient, err := database.NewClient()

		currentTime := time.Now()

		prevDays := currentTime.AddDate(0, 0, -messageBody.Settings.DaysBuffer)

		authorPostCountMap := mongoClient.GetRecentPostCountByAuthorsId(context.TODO(), premiumAuthorsList, prevDays)

		for _, v := range messageBody.Data {
			val, ok := authorPostCountMap[v.UUID]
			if !ok {
				val = 0
			}
			if val < messageBody.Settings.MinPosts {
				nonpremiumAuthorsList = append(nonpremiumAuthorsList, v.UUID)
			}
		}

		// Below Jobs wont run as mongodb creates a copy of current DB and free plan doesnt have much space for it
		// mongoClient.UpdateAuthorPremiumStatus(context.TODO(), premiumAuthorsList, true)
		// mongoClient.UpdateAuthorPremiumStatus(context.TODO(), nonpremiumAuthorsList, false)

		// From below statement we are manually acknowledging the queue so it deletes the message
		d.Ack(false)
	}

	return nil
}
