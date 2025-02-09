package handler

import (
	"url_shorter/shortener"
	"url_shorter/store"

	"net/http"

	"github.com/gin-gonic/gin"
)

// request definition
type UrlCreationRequest struct {
	LongUrl string `json:"long_url" binding:"required"`
	UserId string `json:"user_id" binding:"required"`
}

func CreateShortUrl(c *gin.Context) {
	var creationRequest UrlCreationRequest
	if err := c.ShouldBindJSON(&creationRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shortUrl := shortener.GenertateShortLink(creationRequest.LongUrl,
		creationRequest.UserId)

	store.SaveUrlMapping(shortUrl, creationRequest.LongUrl,
		creationRequest.UserId)

	host := "http://localhost:4440/"
	c.JSON(200, gin.H{
		"message": "shourted url created successfuly",
		"short_url": host+shortUrl,
	})
}

func HandleShourtUrlRedirect(c *gin.Context) {
	shourtUrl := c.Param("shortUrl")
	initialUrl := store.RetrieveInitialUrl(shourtUrl)
	c.Redirect(302, initialUrl)
}
