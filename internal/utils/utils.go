package utils

import (
	"bot/models"
	"log"

	"github.com/gin-gonic/gin"
)

func GetEvent(c *gin.Context) (models.Event, bool) {
	anyEvent, _ := c.Get("event")
	event, ok := anyEvent.(models.Event)
	if !ok {
		log.Println("failed to get event from context")
		return models.Event{}, false
	}
	return event, true
}

func GetApiBaseURL(c *gin.Context) (string, bool) {
	anyBaseURL, _ := c.Get("apiBaseURL")
	apiBaseURL, ok := anyBaseURL.(string)
	if !ok {
		log.Println("failed to get apiBaseURL from context")
		return "", false
	}
	return apiBaseURL, true
}

func GetAuthToken(c *gin.Context) (string, bool) {
	anyAuthToken, _ := c.Get("authToken")
	authToken, ok := anyAuthToken.(string)
	if !ok {
		log.Println("failed to get authToken from context")
		return "", false
	}
	return authToken, true
}
