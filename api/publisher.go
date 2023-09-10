package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"goldeneye.com/m/v2/models"
	"goldeneye.com/m/v2/service"
)

func PostMessage(c *gin.Context) {
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

func UpdateAuthorStatus(c *gin.Context) {
	service.UpdateAuthorStatusJob(c)
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully ran the job",
	})
}

func UpdateGlobalSettings(c *gin.Context) {
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

	service.UpdateGlobalSettings(c, body)
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully updated settings",
	})
}
