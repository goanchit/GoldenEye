package api

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"goldeneye.com/m/v2/repository"
)

func RouteHander(c *gin.Engine, db *mongo.Client) {

	repository := repository.NewClient(db)
	server := NewServer(repository)

	authorRoute := c.Group("/author")
	{
		authorRoute.POST("/post/", server.PostMessage)
		authorRoute.POST("/jobs/", server.UpdateAuthorStatus)
	}
	c.PATCH("/update/settings/", server.UpdateGlobalSettings)

}
