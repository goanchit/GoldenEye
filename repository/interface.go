package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"goldeneye.com/m/v2/models"
)

// All DB related queries
type Repository interface {
	GetUserData(ctx context.Context, document models.MessageBody) *models.User
	GetUserPostCount(ctx context.Context, userId string) int64
	InsertAuthorPost(ctx context.Context, document models.MessageBody) bool
	BulkInsertUserData(ctx context.Context, document []interface{}) error
	GetAllAuthorData(ctx context.Context) ([]models.User, primitive.M)
	GetRecentPostCountByAuthorsId(ctx context.Context, authorIds []string, bufferDate time.Time) map[string]int
	UpdateAuthorPremiumStatus(ctx context.Context, authorIds []string, isPremium bool)
	UpdateAuthorFollowers(ctx context.Context, authorId string, followerCount int)
	UpdateGlobalSettings(ctx context.Context, document interface{})
}
