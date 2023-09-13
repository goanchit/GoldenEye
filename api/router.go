package api

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"goldeneye.com/m/v2/repository"
)

func RouteHander(c *gin.Engine, db *mongo.Client) {

	repository := repository.NewClient(db)
	server := NewServer(repository)

	c.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Server is up and running",
		})
	})

	authorRoute := c.Group("/author")
	{
		authorRoute.POST("/post/", server.PostMessage)
		authorRoute.POST("/jobs/", server.UpdateAuthorStatus)
	}

	consumerRoute := c.Group("/consumer")
	{
		consumerRoute.POST("/post/", server.AuthorPostConsumer)
		consumerRoute.POST("/jobs/", server.DailyAuthorJobConsumer)
	}

	c.PATCH("/update/settings/", server.UpdateGlobalSettings)

}
