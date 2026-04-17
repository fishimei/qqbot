package hook

import (
	"bot/internal/utils"
	"bot/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ParseMsg() gin.HandlerFunc {
	return func(c *gin.Context) {
		var event models.Event
		err := c.BindJSON(&event)
		log.Println(event)
		if err != nil {
			log.Println(err.Error())
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}
		if event.PostType != "message" {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{})
			return
		}
		c.Set("event", event)
	}
}

// JudgeAt 是一个 Gin 中间件函数，用于判断消息事件中是否包含 @ 机器人自己的消息。
func JudgeAt(enable bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !enable {
			return
		}
		event, _ := utils.GetEvent(c)
		if event.MessageType != "group" {
			return
		}
		for _, message := range event.Message {
			if message.Type == "at" && message.Data["qq"] == strconv.FormatInt(event.SelfID, 10) {
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusOK, gin.H{})
	}
}
