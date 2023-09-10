package structures

import "time"

type AuthorPost struct {
	UUID      string    `bson:"uuid"`
	Message   string    `bson:"message"`
	CreatedAt time.Time `bson:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt"`
}

type User struct {
	UUID             string    `bson:"uuid"`
	IsAuthor         bool      `bson:"isAuthor"`
	IsUser           bool      `bson:"isUser"`
	IsPremiumAccount bool      `bson:"isPremiumAccount"`
	CreatedAt        time.Time `bson:"createdAt"`
	UpdatedAt        time.Time `bson:"updatedAt"`
	Followers        int       `bson:"followers,omitempty"`
}
