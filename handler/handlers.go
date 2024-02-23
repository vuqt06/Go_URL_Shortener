package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/vuqt06/go-url-shortener/shortener"
	"github.com/vuqt06/go-url-shortener/store"
)

// Request model definition
type UrlCreationRequest struct {
	LongUrl string `json:"long_url" binding:"required"`
	UserId  string `json:"user_id" binding:"required"`
}

func CreateShortUrl(c *gin.Context) {
	var creationRerquest UrlCreationRequest
	if err := c.ShouldBindJSON(&creationRerquest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shortUrl := shortener.GenerateShortLink(creationRerquest.LongUrl, creationRerquest.UserId)
	store.SaveUrlMapping(shortUrl, creationRerquest.LongUrl, creationRerquest.UserId)

	host := "http://localhost:9808/"
	c.JSON(http.StatusOK, gin.H{
		"short_url": host + shortUrl,
		"message":   "Short URL created successfully",
	})
}

func HandleShortUrlRedirect(c *gin.Context) {
	shortUrl := c.Param("shortUrl")
	longUrl := store.RetrieveOriginalUrl(shortUrl)
	if longUrl == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"})
		return
	}
	c.Redirect(302, longUrl)
}
