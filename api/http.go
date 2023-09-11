package api

import (
	"log"
	"net/http"

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
