package repository

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"goldeneye.com/m/v2/models"
)

type databaseConfig struct {
	client *mongo.Client
}

func NewClient(client *mongo.Client) *databaseConfig {
	return &databaseConfig{
		client: client,
	}
}

func (c databaseConfig) GetUserData(ctx context.Context, document models.MessageBody) *models.User {
	db := c.client.Database(os.Getenv("DATABASE_NAME"))
	var result models.User

	err := db.Collection("User").FindOne(ctx, bson.M{"uuid": document.UserId}).Decode(&result)
	if err != nil {
		log.Fatalf("Error while getting user data %s", err)
		return nil
	}

	return &result
}

func (c databaseConfig) GetUserPostCount(ctx context.Context, userId string) int64 {
	db := c.client.Database(os.Getenv("DATABASE_NAME"))
	authorPostCount, err := db.Collection("AuthorPost").CountDocuments(ctx, bson.M{"uuid": userId})
	if err != nil {
		log.Fatalf("Error while getting user data %s", err)
		return 0
	}

	return authorPostCount
}

// Add new post on database. Increase Post Score by 1
func (c databaseConfig) InsertAuthorPost(ctx context.Context, document models.MessageBody) bool {
	db := c.client.Database(os.Getenv("DATABASE_NAME"))

	err := db.Client().UseSession(ctx, func(sc mongo.SessionContext) error {
		err := sc.StartTransaction()
		if err != nil {
			return err
		}

		user := models.User{
			UUID:             document.UserId,
			IsAuthor:         true,
			IsUser:           false,
			IsPremiumAccount: false,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}

		update := bson.M{
			"$set": user,
		}

		userOptions := options.Update().SetUpsert(true)
		userFilter := bson.M{"uuid": document.UserId}

		// 1. Upsert user if not exists
		_, err = db.Collection("User").UpdateOne(ctx, userFilter, update, userOptions)

		if err != nil {
			log.Fatalf("Failed to create user %s", err)
			sc.AbortTransaction(sc)
			return err
		}

		authorPost := models.AuthorPost{
			UUID:      document.UserId,
			Message:   document.Message,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// 2. Push post to Author Table with reference to uuid
		_, err = db.Collection("AuthorPost").InsertOne(ctx, authorPost)

		if err != nil {
			log.Fatalf("Failed to push author post %s", err)
			sc.AbortTransaction(sc)
			return err
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Insert Transaction Failed %s", err)
	}

	// 3. Get Posts Count and Followers Count Of User
	userData := c.GetUserData(ctx, document)

	return userData.IsPremiumAccount
}

func (c *databaseConfig) BulkInsertUserData(ctx context.Context, document []interface{}) error {
	db := c.client.Database(os.Getenv("DATABASE_NAME"))

	_, err := db.Collection("User").InsertMany(ctx, document)
	if err != nil {
		log.Fatalf("Failed to insert multiple documents %s", err)
	}
	return nil
}

func (c *databaseConfig) GetAllAuthorData(ctx context.Context) ([]models.User, primitive.M) {
	client := c.client
	var result bson.M

	settingsObjectId := os.Getenv("SETTINGS_DOCUMENT_ID")

	objectID, err := primitive.ObjectIDFromHex(settingsObjectId)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Database(os.Getenv("DATABASE_NAME")).Collection("settings").FindOne(ctx, bson.M{"_id": objectID}).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	// Get All Authors Data
	curr, err := client.Database(os.Getenv("DATABASE_NAME")).Collection("User").Find(ctx, bson.M{"isAuthor": true})
	if err != nil {
		log.Println(err)
	}

	var authorDataList []models.User

	for curr.Next(ctx) {
		var authordata models.User
		if err := curr.Decode(&authordata); err != nil {
			log.Printf("Error decoding author data %s", err)
		}
		authorDataList = append(authorDataList, authordata)
	}

	return authorDataList, result
}

func (c *databaseConfig) GetRecentPostCountByAuthorsId(ctx context.Context, authorIds []string, bufferDate time.Time) map[string]int {
	db := c.client.Database(os.Getenv("DATABASE_NAME"))

	m := make(map[string]int)

	for _, id := range authorIds {
		count, err := db.Collection("AuthorPost").CountDocuments(ctx, bson.M{"uuid": id, "createdAt": bson.M{"$gte": primitive.NewDateTimeFromTime(bufferDate)}})
		if err != nil {
			log.Fatal("Failed to get count")
		}
		m[id] = int(count)
	}
	return m
}

func (c *databaseConfig) UpdateAuthorPremiumStatus(ctx context.Context, authorIds []string, isPremium bool) {
	db := c.client.Database(os.Getenv("DATABASE_NAME"))

	_, err := db.Collection("User").UpdateMany(ctx, bson.M{"uuid": bson.M{"$in": authorIds}}, bson.M{"$set": bson.M{"isPremiumAccount": isPremium}})
	if err != nil {
		log.Fatal(err)
		log.Fatalf("Failed to update author status for %t", isPremium)
	}
	return
}

func (c *databaseConfig) UpdateAuthorFollowers(ctx context.Context, authorId string, followerCount int) {
	db := c.client.Database(os.Getenv("DATABASE_NAME"))
	_, err := db.Collection("User").UpdateOne(ctx, bson.M{"uuid": authorId}, bson.D{{"$inc", bson.D{{"followers", followerCount}}}})
	if err != nil {
		log.Fatal(err)
		log.Fatal("Failed to update author followers")
	}
	return
}

func (c *databaseConfig) UpdateGlobalSettings(ctx context.Context, document interface{}) {
	db := c.client.Database(os.Getenv("DATABASE_NAME"))

	settingsObjectId := os.Getenv("SETTINGS_DOCUMENT_ID")
	objectID, err := primitive.ObjectIDFromHex(settingsObjectId)

	if err != nil {
		log.Fatal(err)
	}

	result, err := db.Collection("settings").UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": document})
	if err != nil {
		log.Fatal(err)
		log.Fatal("Failed to update settings")
	}

	log.Printf("Matched %v documents and modified %v documents\n", result.MatchedCount, result.ModifiedCount)
	return
}
