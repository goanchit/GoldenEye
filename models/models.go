package models

import "goldeneye.com/m/v2/structures"

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
	Data     []structures.User `json:"data"`
	Settings Settings          `json:"settings"`
}
