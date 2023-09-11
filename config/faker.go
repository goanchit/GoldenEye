package config

import (
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

type FakeUserStruct struct {
	UUID             string    `bson:"uuid"`
	IsAuthor         bool      `bson:"isAuthor"`
	IsUser           bool      `bson:"isUser"`
	IsPremiumAccount bool      `bson:"isPremiumAccount"`
	CreatedAt        time.Time `bson:"createdAt"`
	UpdatedAt        time.Time `bson:"updatedAt"`
	Followers        int       `bson:"followers"`
}

func convertToInterfaceSlice(inputSlice []FakeUserStruct) []interface{} {
	var outputSlice []interface{}

	for _, item := range inputSlice {
		outputSlice = append(outputSlice, item)
	}

	return outputSlice
}

// Generate Fake User Data which will depict users
func GenerateFakeData() []interface{} {
	var fakeData []FakeUserStruct
	i := 1

	for i < 200 {
		var res FakeUserStruct
		res.CreatedAt = gofakeit.Date()
		res.UpdatedAt = res.CreatedAt
		res.Followers = 0
		res.IsAuthor = false
		res.IsPremiumAccount = gofakeit.Bool()
		res.IsUser = true
		res.UUID = gofakeit.UUID()

		fakeData = append(fakeData, res)
		i++
	}
	// Converts to Type []interface{} this is required while bulk insert
	res := convertToInterfaceSlice(fakeData)

	return res
}
