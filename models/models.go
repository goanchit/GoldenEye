package models

import (
	"time"
)

type MessageBody struct {
	UserId  string `json:"userId" validate:"required,uuid"`
	Message string `json:"message" validate:"required,min=10,max=200"`
}

// map[days_buffer:7 min_followers:20 min_posts:2]}
type Settings struct {
	DaysBuffer   int `json:"days_buffer" bson:"days_buffer,omitempty"`
	MinFollowers int `json:"min_followers" bson:"min_followers,omitempty"`
	MinPosts     int `json:"min_posts" bson:"min_posts,omitempty"`
}

type AuthorPremiumJob struct {
	Data     []User   `json:"data"`
	Settings Settings `json:"settings"`
}

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
