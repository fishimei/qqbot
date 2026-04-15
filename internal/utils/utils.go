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
