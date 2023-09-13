package api

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"goldeneye.com/m/v2/models"
	"goldeneye.com/m/v2/repository"
	"goldeneye.com/m/v2/service"
)

type Server struct {
	repository repository.Repository
}

func NewServer(repository repository.Repository) *Server {
	return &Server{
		repository: repository,
	}
}

func (s Server) PostMessage(c *gin.Context) {
	body := models.MessageBody{}

	if err := c.ShouldBind(&body); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad Body Parsed",
		})
		return
	}

	validate := validator.New()

	if err := validate.Struct(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.(validator.ValidationErrors),
		})
		return
	}

	service.PublishToQueue(c, "AUTHOR_POST", body)

	c.JSON(http.StatusOK, gin.H{
		"message": "Post Published",
	})
}

func (s Server) UpdateAuthorStatus(c *gin.Context) {

	allAuthors, settings := s.repository.GetAllAuthorData(c)

	batches := chunkBy(allAuthors, 4)

	for _, v := range batches {
		m := make(map[string]interface{})

		m["data"] = v
		m["settings"] = settings

		service.PublishToQueue(c, "AUTHOR_STATUS_JOB", m)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully ran the job",
	})
}

func (s Server) UpdateGlobalSettings(c *gin.Context) {
	body := models.Settings{}
	if err := c.ShouldBind(&body); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad Body Parsed",
		})
		return
	}

	validate := validator.New()

	if err := validate.Struct(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.(validator.ValidationErrors),
		})
		return
	}

	s.repository.UpdateGlobalSettings(c, body)

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully updated settings",
	})
}

// Convert Array to Slices of size chunkSize
func chunkBy[T any](items []T, chunkSize int) (chunks [][]T) {
	for chunkSize < len(items) {
		items, chunks = items[chunkSize:], append(chunks, items[0:chunkSize:chunkSize])
	}
	return append(chunks, items)
}

func (s Server) AuthorPostConsumer(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)

	if err != nil {
		log.Fatal("Error reading the body: ", err)
		return
	}

	byteString, err := base64.StdEncoding.DecodeString(string(body))
	if err != nil {
		log.Fatal("Error failed to Decode Message Body: ", err)
		return
	}

	var messageBody models.MessageBody
	err = json.Unmarshal(byteString, &messageBody)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Message Body Recieved: ", messageBody)

	isPremiumAuthor := s.repository.InsertAuthorPost(c, messageBody)

	log.Print("Is Premium Author: ", isPremiumAuthor)

	// Complete this
	if isPremiumAuthor {
		log.Print("Sending Post To Premium User")
		// Send Post to premium users
	} else {
		log.Print("Sending Post To Regular User")
		// Send Post to regular users
	}

	// Below code depicts number of random followers to be added
	randomFollowers := rand.Intn(10)
	s.repository.UpdateAuthorFollowers(c, messageBody.UserId, randomFollowers)

	return
}

func (s Server) DailyAuthorJobConsumer(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)

	if err != nil {
		log.Fatal("Error reading the body: ", err)
		return
	}

	byteString, err := base64.StdEncoding.DecodeString(string(body))
	if err != nil {
		log.Fatal("Error failed to Decode Message Body: ", err)
		return
	}

	var messageBody models.AuthorPremiumJob
	err = json.Unmarshal(byteString, &messageBody)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Message Body Recieved: ", messageBody)

	// For Premium Authors set premium flag to true
	// For NonPremium Authors set premium flag to false
	var premiumAuthorsList, nonpremiumAuthorsList []string

	for _, v := range messageBody.Data {
		if messageBody.Settings.MinFollowers <= v.Followers {
			premiumAuthorsList = append(premiumAuthorsList, v.UUID)
		} else {
			nonpremiumAuthorsList = append(nonpremiumAuthorsList, v.UUID)
		}
	}

	currentTime := time.Now()
	prevDays := currentTime.AddDate(0, 0, -messageBody.Settings.DaysBuffer)

	// Pull this from inmemory
	authorPostCountMap := s.repository.GetRecentPostCountByAuthorsId(c, premiumAuthorsList, prevDays)

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
	// mongoClient.UpdateAuthorPremiumStatus(ctx, premiumAuthorsList, true)
	// mongoClient.UpdateAuthorPremiumStatus(ctx, nonpremiumAuthorsList, false)

	return
}
