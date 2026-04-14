package hook

import (
	"bot/internal/utils"
	"bot/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func PassToken(apiBaseURL, authToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("apiBaseURL", apiBaseURL)
		c.Set("authToken", authToken)
		if authToken == "" {
			c.AbortWithStatusJSON(http.StatusExpectationFailed, gin.H{"error": "客户端authToken未配置，请检查配置文件"})
		}
	}
}

// CheckToken 是一个 Gin 中间件函数，用于验证请求中的 Authorization 头部是否包含预期的 Bearer Token。
func CheckToken(expectedToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if expectedToken != "" {
			token := c.GetHeader("Authorization")
			if token != "Bearer "+expectedToken {
				log.Println(token)
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token is invalid"})
			}
		}
	}
}

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

// WhiteList todo
func WhiteList() gin.HandlerFunc {
	return func(c *gin.Context) {

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
