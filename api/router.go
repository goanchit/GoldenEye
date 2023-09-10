package api

import (
	"github.com/gin-gonic/gin"
)

func RouteHander(c *gin.Engine) {

	authorRoute := c.Group("/author")
	{
		authorRoute.POST("/post/", PostMessage)
		authorRoute.POST("/jobs/", UpdateAuthorStatus)
	}
	c.PATCH("/update/settings/", UpdateGlobalSettings)

}
